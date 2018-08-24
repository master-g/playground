package agent

import (
	"encoding/binary"
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrorSendEmptyPacket error while sending empty packet
	ErrorSendEmptyPacket = errors.New("sending invalid packet (nil)")
	// ErrorQueueFull error while tx queue is full
	ErrorQueueFull = errors.New("sending queue is full")
)

// ConnWriteFunc connection writer function type
type ConnWriteFunc func([]byte) (int, error)

// Sender send packets to the client
type Sender struct {
	ctrl    chan struct{} // exit signal
	pending chan []byte   // pending packets
	writer  ConnWriteFunc // connection writer function
	cache   []byte        // for combined syscall write
}

// NewSender create and returns a sender instance
func NewSender(ctrl chan struct{}, queueSize, sendCacheSize int, writer ConnWriteFunc) *Sender {
	return &Sender{
		ctrl:    ctrl,
		writer:  writer,
		pending: make(chan []byte, queueSize),
		cache:   make([]byte, sendCacheSize),
	}
}

// EnqueueOutgoing push data to the sender's pending channel
func (buf *Sender) EnqueueOutgoing(pkg []byte) error {
	if pkg == nil {
		return ErrorSendEmptyPacket
	}

	// queue the data for sending
	select {
	case buf.pending <- pkg:
	default:
		log.Warn(ErrorQueueFull.Error())
		return ErrorQueueFull
	}

	return nil
}

// actual send logic
func (buf *Sender) actualSend(data []byte) bool {
	// write packet size, uint16, (these 2 bytes are excluded)
	size := len(data)
	binary.BigEndian.PutUint16(buf.cache, uint16(size))
	copy(buf.cache[2:], data)

	// write data
	n, err := buf.writer(buf.cache[:size+2])
	if err != nil {
		log.Errorf("error while sending data %v bytes, %v", n, err)
		return false
	}

	return true
}

// SendLoop packets sending routine
func (buf *Sender) SendLoop() {
	for {
		select {
		case data := <-buf.pending: // dequeue data for sending
			buf.actualSend(data)
		case <-buf.ctrl: // control signal received
			// TODO
			// 1. how can we send last packet before close
			// 2. only send reason packet before close
			return
		}
	}
}
