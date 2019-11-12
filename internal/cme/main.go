package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func fetch(repo string) error {
	d := time.Duration(rand.Intn(5)) * time.Second
	<-time.After(d)
	if rand.Intn(100) > 95 {
		return fmt.Errorf("rua (%v)", repo)
	}

	return nil
}

func restore(repos []string) error {
	errChan := make(chan error, 1)
	sem := make(chan int, 4)
	var wg sync.WaitGroup
	wg.Add(len(repos))

	for _, repo := range repos {
		sem <- 1
		go func() {
			defer func() {
				wg.Done()
				<-sem
			}()
			if err := fetch(repo); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(sem)
	close(errChan)
	return <-errChan
}

func main() {

}
