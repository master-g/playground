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
	// FlagAuth indicates the session has been authorized
	FlagAuth = 0x8
)

// Session holds the context of a client having conversation with agent
type Session struct {
	Die               chan struct{}
	flag              Flag
	IP                net.IP
	Port              string
	MQ                chan []byte
	Push              chan []byte
	ConnectTime       time.Time
	PacketTime        time.Time
	LastPacketTime    time.Time
	PacketCount       int
	PacketCountPerMin int
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

// SetFlagAuth sets the auth bit
func (s *Session) SetFlagAuth() *Session {
	s.flag |= FlagAuth
	return s
}

// ClearFlagAuth clears the auth bit
func (s *Session) ClearFlagAuth() *Session {
	s.flag &^= FlagAuth
	return s
}

// IsFlagAuthSet returns true if the auth bit is set
func (s *Session) IsFlagAuthSet() bool {
	return s.flag&FlagAuth != 0
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

// CheckRPMLimitViolation returns true if session violates RPM limitation
func (s *Session) CheckRPMLimitViolation() bool {
	if s.PacketCountPerMin > s.RPMLimit {
		log.Infof("RPM violation, session: %v, rate: %v, total: %v", s.String(), s.PacketCountPerMin, s.PacketCount)
		return true
	} else {
		return false
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
