package agent

import (
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/lopnur/lnutils/signal"
	log "github.com/sirupsen/logrus"
)

// StartTCPClient dials server tcp and send messages
func StartTCPClient() {
	addr, err := net.ResolveTCPAddr("tcp4", ":8888")
	if err != nil {
		log.Fatalf("unable to resolve TCP address, error %v", err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("unable to dial TCP address, error %v", err)
	}

	go clientLoop(conn)
}

func clientLoop(conn net.Conn) {
	defer conn.Close()

	packetSize := make([]byte, 2)
	in := make(chan []byte)
	defer func() {
		close(in)
	}()

	// send loop
	ctrl := make(chan struct{})
	out := NewSender(ctrl, 128, 32768, func(data []byte) (n int, err error) {
		return conn.Write(data)
	})
	go out.SendLoop()

	// handling incoming packet
	go handleIncomingPacket(ctrl, in)

	// send heartbeat packet to server
	go tick(out)

	// read loop
	for {
		n, err := io.ReadFull(conn, packetSize)
		if err != nil {
			log.Errorf("read packet size failed %v bytes read, error %v", n, err)
			signal.RequestInterrupt()
			return
		}
		size := binary.BigEndian.Uint16(packetSize)

		// read packet
		packetData := make([]byte, size)
		n, err = io.ReadFull(conn, packetData)
		if err != nil {
			log.Errorf("read packet body failed %v expected, %v read, error %v", size, n, err)
			signal.RequestInterrupt()
			return
		}

		// deliver
		select {
		case in <- packetData:
		case <-ctrl:
			log.Info("connection closed by logic")
			return
		}
	}
}

func handleIncomingPacket(ctrl chan struct{}, in chan []byte) {
	defer func() {
		close(ctrl)
	}()

	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}
			log.Infof("[RECV] %v", string(msg))
		}
	}
}

func tick(out *Sender) {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			out.EnqueueOutgoing([]byte(time.Now().String()))
		}
	}
}
