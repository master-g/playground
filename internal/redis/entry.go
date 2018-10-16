package redis

import (
	"fmt"
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func Entry() {
	redis := &Redis{}
	redis.Init(&Config{
		Network:     "tcp",
		Address:     ":6379",
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: time.Second * 5,
		Wait:        true,
	})

	d, err := redis.Get("shit")
	if err != nil {
		log.Fatal(err)
	}
	if d == nil {
		log.Println("value is nil")
	}
	s, err := redigo.String(d, err)
	fmt.Println(s)
}
