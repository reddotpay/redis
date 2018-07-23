package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// Pool contains Redis pool
var Pool *redis.Pool

// NewPool initialises a new Redis pool
func NewPool(endpoint string, maxIdle int, idleTimeout time.Duration) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(endpoint)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

// Store stores key-value pairs in Redis with expiry
func Store(key string, value interface{}, expiryInSeconds int) error {
	conn := Pool.Get()

	conn.Send("MULTI")
	conn.Send("SET", key, value)

	if expiryInSeconds > 0 {
		conn.Send("EXPIRE", key, expiryInSeconds)
	}

	_, err := conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

// Retrieve retrieves value by key
func Retrieve(key string) (interface{}, error) {
	conn := Pool.Get()
	reply, err := conn.Do("GET", key)
	if err != nil {
		if err == redis.ErrNil {
			return reply, nil
		}

		return reply, err
	}

	return reply, err
}
