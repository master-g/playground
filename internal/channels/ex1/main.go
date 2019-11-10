package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

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

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}
