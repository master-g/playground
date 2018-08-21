// Copyright © 2018 Project Lop Nur <project.lopnur@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Network     string        // tcp
	Address     string        // ip:port
	Password    string        // password
	MaxIdle     int           // maximum number of idle connections
	MaxActive   int           // maximum number of active connections
	IdleTimeout time.Duration // idle connection timeout
	Wait        bool          // wait for connection returned to pool before get
}

// Redis struct
type Redis struct {
	pool   *redis.Pool // redis connection pool
	config *Config
}

// Init redis with parameters
func (r *Redis) Init(cfg *Config) {
	r.config = cfg
}

func dial(network, address, password string) (c redis.Conn, err error) {
	c, err = redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err = c.Do("AUTH", password); err != nil {
			c.Close()
			c = nil
			return
		}
	}
	return
}

// GetPool get redis pool
func (r *Redis) GetPool() *redis.Pool {
	if r.pool == nil {
		r.pool = &redis.Pool{
			MaxIdle:     r.config.MaxIdle,
			MaxActive:   r.config.MaxActive,
			Wait:        r.config.Wait,
			IdleTimeout: r.config.IdleTimeout,
			Dial: func() (redis.Conn, error) {
				return dial(r.config.Network, r.config.Address, r.config.Password)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < 10*time.Second {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}
	}
	return r.pool
}

// Put a key-value pair to redis with a timeout
func (r *Redis) Put(key string, value interface{}, time int) (interface{}, error) {
	conn := r.GetPool().Get()
	defer conn.Close()

	reply, err := conn.Do("SETEX", key, time, value)
	// reply, err := redis.String(conn.Do("SETEX", key, time, value))

	return reply, err
}

// Set a key-value pair to redis
func (r *Redis) Set(key string, value interface{}) (interface{}, error) {
	conn := r.GetPool().Get()
	defer conn.Close()

	reply, err := conn.Do("SET", key, value)

	return reply, err
}

// Get a key-value pair from redis
func (r *Redis) Get(key string) (interface{}, error) {
	conn := r.GetPool().Get()
	defer conn.Close()

	reply, err := conn.Do("GET", key)

	return reply, err
}

// Del a key-value pair from redis
func (r *Redis) Del(key string) (interface{}, error) {
	conn := r.GetPool().Get()
	defer conn.Close()

	reply, err := conn.Do("DEL", key)

	return reply, err
}

// Close redis connection pool
func (r *Redis) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}
