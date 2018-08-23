package socket

import (
	"log"
	"net"
	"time"
)

func Entry() {
	conn, err := net.Dial("tcp", ":8888")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 30)
}
