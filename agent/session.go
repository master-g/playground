package agent

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Flag holds session status flag bits
type Flag int32

const (
	// FlagKicked indicates the client has been kicked out
	FlagKicked = 0x4
)

// Session holds the context of a client having conversation with agent
type Session struct {
	Die               chan struct{}
	flag              Flag
	IP                net.IP
	Port              string
	MQ                chan []byte
	ConnectTime       time.Time
	PacketTime        time.Time
	LastPacketTime    time.Time
	PacketCount       uint32
	PacketCountPerMin int
	Push              chan []byte
	RPMLimit          int
}

// NewSession returns a new instance of session
func NewSession(ip net.IP, port string, rpmLimit int) *Session {
	s := &Session{
		IP:       ip,
		Port:     port,
		Die:      make(chan struct{}),
		RPMLimit: rpmLimit,
	}

	return s
}

// SetFlagKicked sets the kicked bit
func (s *Session) SetFlagKicked() *Session {
	s.flag |= FlagKicked
	return s
}

// ClearFlagKicked clears the kicked bit
func (s *Session) ClearFlagKicked() *Session {
	s.flag &^= FlagKicked
	return s
}

// IsFlagKickedSet returns true if the kicked bit is set
func (s *Session) IsFlagKickedSet() bool {
	return s.flag&FlagKicked != 0
}

// TimeWork checks rpm limit
func (s *Session) TimeWork() {
	defer func() {
		s.PacketCountPerMin = 0
	}()

	// rpm control
	if s.PacketCountPerMin > s.RPMLimit {
		s.SetFlagKicked()
		log.Errorf("RPM violation %v", s.String())
		return
	}
}

// FetchLoop fetches streams from game service
func (s *Session) FetchLoop() {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			s.MQ <- []byte(fmt.Sprintf("stream from game server %v", time.Now()))
		case <-s.Die:
			log.Infof("session closed by logic, stop streaming %v", s.String())
			return
		}
	}
}

// String interface
func (s *Session) String() string {
	return fmt.Sprintf("ip:%v port:%v flag:%v  ppm:%v limit:%v", s.IP.String(), s.Port, s.flag, s.PacketCountPerMin, s.RPMLimit)
}
