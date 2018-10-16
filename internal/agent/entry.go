package agent

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lopnur/lnutils/signal"
	log "github.com/sirupsen/logrus"
)

func printUsage() {
	fmt.Println("usage: playground [client|server]")
	os.Exit(0)
}

// Entry for agent module
func Entry() {
	serverFlag := flag.Bool("server", false, "run in server mode")
	clientFlag := flag.Bool("client", false, "run in client mode")

	flag.Parse()

	if *serverFlag == *clientFlag {
		printUsage()
	}

	if *serverFlag {
		startServer()
		log.Info("bye")
	} else {
		startClient()
		log.Info("bye")
	}
}

func startServer() {
	go StartTCP(&Config{
		Port:             8888,
		ReadDeadLine:     60 * time.Second,
		ReadBufferSize:   32768,
		WriteDeadLine:    10 * time.Second,
		WriteBufferSize:  32768,
		SessionCacheSize: 32768,
		TxQueueLength:    128,
		RPMLimit:         200,
		AuthTimeout:      1 * time.Second,
		CloseTimeout:     1 * time.Second,
	})

	go signal.Start()
	<-signal.InterruptChan
	WaitTCPShutdown()
}

func startClient() {
	go StartTCPClient()
	go signal.Start()
	<-signal.InterruptChan
}
