package wg

import (
	"math/rand"
	"sync"
	"time"
)

func task(wg *sync.WaitGroup) {
	defer wg.Done()
	waitDuration := time.Duration(rand.Int63n(3)) * time.Second
	timer := time.NewTimer(waitDuration)
	<-timer.C
}

func Execute() {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go task(&wg)
	}
	wg.Wait()
}
