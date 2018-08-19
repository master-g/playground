package signal

import (
	"os"
	"os/signal"
	"syscall"
)

type InterruptType int

const (
	SystemInterrupt InterruptType = 0
	UserInterrupt   InterruptType = 1
)

var (
	// InterruptChan as signal to control goroutine to finish
	InterruptChan = make(chan InterruptType)
	// interruptRequest for other subsystem to initiate a proper shutdown
	interruptRequest = make(chan struct{})
	// interruptSignals defines system signals to catch
	interruptSignals = []os.Signal{syscall.SIGTERM, os.Interrupt}
)

// Start a goroutine to process UNIX system signal
func Start() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, interruptSignals...)

	select {
	case <-c:
		InterruptChan <- SystemInterrupt
	case <-interruptRequest:
		InterruptChan <- UserInterrupt
	}

	close(InterruptChan)
}

// RequestInterrupt from other subsystem to initiate a proper shutdown
func RequestInterrupt() {
	close(interruptRequest)
}
