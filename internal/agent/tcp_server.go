package agent

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/lopnur/lnutils/signal"
	log "github.com/sirupsen/logrus"
)

// Config tcp server
type Config struct {
	Port             int
	ReadDeadLine     time.Duration
	ReadBufferSize   int
	WriteDeadLine    time.Duration
	WriteBufferSize  int
	SessionCacheSize int
	TxQueueLength    int
	RPMLimit         int
	AuthTimeout      time.Duration
	CloseTimeout     time.Duration
}

var (
	wg     *sync.WaitGroup
	config *Config
)

func init() {
	wg = &sync.WaitGroup{}
}

// WaitTCPShutdown waits all connection to close
func WaitTCPShutdown() {
	wg.Wait()
}

// StartTCP start TCP server
func StartTCP(cfg *Config) {
	defer signal.RequestInterrupt()
	if cfg == nil {
		log.Error("nil config")
		return
	}
	config = cfg

	// resolve host address
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Errorf("unable to resolve TCP address %v", err.Error())
		return
	}

	// start TCP listening
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("unable to start listen TCP %v", err.Error())
		return
	}

	log.Infof("listening %v", listener.Addr().String())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Errorf("unable to accept connection %v", err.Error())
			continue
		}

		// setup socket
		conn.SetReadBuffer(config.ReadBufferSize)
		conn.SetWriteBuffer(config.WriteBufferSize)
		// start a goroutine for every incoming connection
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	// close connection when all is done
	defer conn.Close()

	// agent's input channel
	in := make(chan []byte)
	defer func() {
		close(in)
	}()

	// create a new session object for this connection
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Errorf("unable to parse remote address %v", err.Error())
		return
	}

	sess := NewSession(net.ParseIP(host), port, config.RPMLimit)
	go sess.FetchLoop()
	log.Infof("new connection %v", sess.String())

	// create sender and start sending loop
	out := NewSender(sess.Die, config.TxQueueLength, config.SessionCacheSize, func(data []byte) (n int, err error) {
		return conn.Write(data)
	})
	go out.SendLoop()

	// start agent for packet processing
	wg.Add(1)
	go agent(wg, sess, in, out)

	// packet size
	packetSize := make([]byte, 2)

	// killer
	if config.CloseTimeout != 0 {
		go func() {
			<-sess.Die
			time.Sleep(config.CloseTimeout)
			conn.Close()
		}()
	}

	for {
		// solve dead link problem
		conn.SetReadDeadline(time.Now().Add(config.ReadDeadLine))

		// read packet size
		n, err := io.ReadFull(conn, packetSize)
		if err != nil {
			log.Errorf("read packet size failed, error %v session %v", err, sess.String())
			return
		}
		size := binary.BigEndian.Uint16(packetSize)

		// read packet body
		packetData := make([]byte, size)
		n, err = io.ReadFull(conn, packetData)
		if err != nil {
			log.Errorf("read packet body failed, expect %v actual %v session %v", size, n, sess.String())
			return
		}

		select {
		case in <- packetData:
			// deliver the payload to the input queue of agent
		case <-sess.Die:
			// wait for session close in agent
			log.Infof("connection closed by logic, stop reading from connection, session %v", sess.String())
			return
		}
	}
}
