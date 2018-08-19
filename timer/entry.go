package timer

import (
	"fmt"
	"time"

	"github.com/lopnur/lnutils/signal"
)

func Entry() {
	go signal.Start()
	tick()
	fmt.Println("graceful exit")
}

func tick() {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()

	go func() {
		for ; true; <-t.C {
			fmt.Println("tick")
		}
	}()

	<-signal.InterruptChan
}
