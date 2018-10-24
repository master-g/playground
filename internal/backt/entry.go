package backt

import (
	"errors"
	"math/rand"

	"github.com/cenkalti/backoff"
	"github.com/lopnur/kroran/pkg/signal"
	log "github.com/sirupsen/logrus"
)

// Entry for this package
func Entry() {
	go signal.Start()
	go worker()

	<-signal.InterruptChan
	log.Info("bye~")
}

func worker() {
	b := backoff.NewExponentialBackOff()
	ticker := backoff.NewTicker(b)
	defer ticker.Stop()
	var err error

	for {
		select {
		case <-ticker.C:
			err = operation()
			if err != nil {
				log.Errorf("synchronize error %v", err)
				i := b.NextBackOff()
				log.Info(i)
			} else {
				log.Info("synchronize finished")
				b.Reset()
			}
		case <-signal.InterruptChan:
			return
		}
	}
}

func operation() error {
	log.Info("op...")
	if rand.Intn(10) < 5 {
		return errors.New("wow")
	}
	return nil
}
