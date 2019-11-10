package timer

import (
	"fmt"
	"time"

	"github.com/master-g/playground/pkg/signal"
)

func Entry() {
	go signal.Start()
	// go tick()
	c := make(chan int)
	go producer(c)
	go consumer(c)
}

func tick() {
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			fmt.Println("tick")
		case <-signal.InterruptChan:
			fmt.Println("good bye")
			return
		}
	}
}

func consumer(s chan int) {
	consumeTicker := time.NewTicker(time.Second * 30)
	defer consumeTicker.Stop()
	for {
		select {
		case <-consumeTicker.C:
			select {
			case <-s:
				fmt.Println("got from channel")
			}
		case <-signal.InterruptChan:
			return
		}
	}
}

func producer(s chan int) {
	pushTicker := time.NewTicker(time.Second * 5)
	defer pushTicker.Stop()

	for {
		select {
		case <-pushTicker.C:
			select {
			case s <- 233:
			default:
				fmt.Println("s is blocking")
			}

			// select {
			// case s <- 233:
			// 	fmt.Println("write to s")
			// }
		case <-signal.InterruptChan:
			fmt.Println("producer goodbye")
			return
		default:
			time.Sleep(time.Second)
			fmt.Println(time.Now())
		}
	}
}
