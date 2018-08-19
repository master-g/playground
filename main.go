package main

import (
	"fmt"

	"github.com/lopnur/lnutils/signal"
	"github.com/master-g/playground/timer"
)

func main() {
	go timer.Entry()
	fmt.Println("wow")
	<-signal.InterruptChan
}
