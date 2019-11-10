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
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	dataCh := make(chan int)
	stopCh := make(chan struct{})

	for i := 0; i < NumSenders; i++ {
		go func() {
			for {
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- rand.Intn(Max):
				}
			}
		}()
	}

	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			if value == Max-1 {
				close(stopCh)
				return
			}

			log.Println(value)
		}
	}()

	wgReceivers.Wait()
}
