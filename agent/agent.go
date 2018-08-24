package agent

import (
	"sync"
	"time"

	"github.com/lopnur/lnutils/signal"
	log "github.com/sirupsen/logrus"
)

const (
	defaultMQSize        = 512
	defaultPushQueueSize = 128
)

// all packet from handleTCPConnection() will be handle here
func agent(wg *sync.WaitGroup, s *Session, in chan []byte, out *Sender) {
	defer wg.Done()

	// init session
	s.MQ = make(chan []byte, defaultMQSize)
	s.Push = make(chan []byte, defaultPushQueueSize)
	s.ConnectTime = time.Now()
	s.LastPacketTime = time.Now()
	// minute timer
	rpmLimitTicker := time.NewTicker(time.Minute)

	// cleanup
	defer func() {
		// notify handleTCPConnection()
		close(s.Die)
	}()

	// MAIN MESSAGE LOOP
	for {
		select {
		case msg, ok := <-in:
			// process packet from network
			if !ok {
				log.Infof("incoming packet channel full/or closed, session %v", s.String())
				s.SetFlagKicked()
				break
			} else {
				// route packet
				s.PacketCount++
				s.PacketCountPerMin++
				s.PacketTime = time.Now()

				// s.SetFlagKicked()

				// route
				echo := []byte("echo ")
				response := append(echo, msg...)
				sendPacket(s, out, response)
				s.LastPacketTime = time.Now()
			}
		case msg := <-s.Push:
			// internal push
			sendPacket(s, out, msg)
		case frame := <-s.MQ:
			sendPacket(s, out, frame)
		case <-rpmLimitTicker.C:
			s.TimeWork()
		case <-signal.InterruptChan:
			// server is manually shutting down
			s.SetFlagKicked()
		}

		if s.IsFlagKickedSet() {
			log.Infof("session kicked %v", s.String())
			return
		}
	}
}

func sendPacket(s *Session, buf *Sender, pkg []byte) {
	err := buf.EnqueueOutgoing(pkg)
	if err != nil {
		log.Errorf("error while sending, session %v error %v", s.String(), err)
	}
}
