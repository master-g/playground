package concu

import (
	"sync"

	"github.com/lopnur/lnutils/signal"
	"github.com/onrik/ethrpc"
	log "github.com/sirupsen/logrus"
)

type CC struct {
	sync.RWMutex
	client  *ethrpc.EthRPC
	version string
}

// Entry for this package
func Entry() {
	client := ethrpc.New("http://localhost:8545")
	cc := &CC{
		client: client,
	}

	for i := 0; i < 999; i++ {
		// go getChainID(client)
		go cc.Version()
	}

	<-signal.InterruptChan
}

func (c *CC) Version() string {
	c.RLock()

	if c.version != "" {
		c.RUnlock()
		log.Info(c.version)
		return c.version
	}

	c.RUnlock()
	c.Get()

	log.Info(c.version)

	return c.version
}

func (c *CC) Get() {
	if c.version != "" {
		return
	}

	c.Lock()
	defer c.Unlock()

	ver, err := c.client.NetVersion()
	if err != nil {
		log.Error(err)
		return
	}

	c.version = ver
}
