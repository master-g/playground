package main

import (
	"flag"
	"math/rand"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func premute() {
	poolorder := make([]int, 8)
	// loop through slice
	for i := 1; i < len(poolorder); i++ {
		j := rand.Int31n(int32(i + 1)) // choose a random index between 0~i
		poolorder[i] = poolorder[j]    // swap slice[i] and slice[j]
		poolorder[j] = i               // slice[j] = i
	}

	log.Info(poolorder)
}

func fileServer() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "directory of static files to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Infof("Serving %v on HTTP port: %v", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func main() {
	fileServer()
}
