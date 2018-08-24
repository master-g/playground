package agent

import (
	"fmt"
	"os"
	"strings"
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
	if len(os.Args) < 2 {
		printUsage()
	}
	if strings.ToLower(os.Args[1]) == "server" {
		startServer()
		log.Info("bye")
	} else if strings.ToLower(os.Args[1]) == "client" {
		startClient()
		log.Info("bye")
	} else {
		printUsage()
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
