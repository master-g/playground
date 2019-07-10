package main

import (
	"fmt"
	"time"
)

var (
	stopCh chan struct{}
)

func foo() {
	tick := time.Tick(time.Second)
LOOP:
	for {
		select {
		case <-stopCh:
			fmt.Println("stop signal received at first select")
			break LOOP
		default:
		}

		select {
		case <-stopCh:
			fmt.Println("stop signal received at second select")
			break LOOP
		case <-tick:
			fmt.Print("some heavy job...")
			<-time.After(500 * time.Millisecond)
			fmt.Println("done")
		}
	}
}

func main() {
	stopCh = make(chan struct{})
	go foo()

	<-time.After(5 * time.Second)
	close(stopCh)

	select {}
}
