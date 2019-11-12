package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// M receivers, one sender

func main() {
	rand.Seed(time.Now().Unix())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 100

	wgReceivers := sync.WaitGroup{}

	dataCh := make(chan int)

	// sender
	go func() {
		for {
			if value := rand.Intn(Max); value == 0 {
				// The only sender can close the
				// channel at any time safely.
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	// receivers
	for i := 0; i < NumReceivers; i++ {
		wgReceivers.Add(1)
		go func() {
			defer wgReceivers.Done()

			// Receive values until dataCh is
			// closed and the value buffer queue
			// of dataCh becomes empty.
			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}
