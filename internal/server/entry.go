package server

import (
	"fmt"
	"time"

	"github.com/master-g/playground/pkg/signal"
)

var (
	die chan struct{}
	mq  chan string
)

// Entry point of this package
func Entry() {
	die = make(chan struct{})
	mq = make(chan string, 1)

	defer func() {
		close(die)
		close(mq)
	}()

	// simulate agent
	go agent()

	// start dummy server
	go server("0")
	go server("1")

	// close sesison after a period of time
	go func() {
		time.Sleep(time.Second * 60)
		close(die)
	}()

	go signal.Start()

	// wait for shutdown
	<-signal.InterruptChan
	fmt.Println("gracefully shuting down")
	// graceful shutdown
	time.Sleep(time.Second)
	fmt.Println("bye")
}

func agent() {
	for {
		select {
		case msg := <-mq:
			// network latency
			fmt.Println(msg)
		case <-die:
			fmt.Println("session die in agent")
			return
		}
	}
}

func server(id string) {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			mq <- fmt.Sprintf("game msg %v", id)
		case <-signal.InterruptChan:
			return
		}
	}
}
