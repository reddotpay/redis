package redis

import (
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

// Pool contains Redis pool
var Pool *redigo.Pool

// NewPool initialises a new Redis pool
func NewPool(endpoint string, maxIdle int, idleTimeout time.Duration) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(endpoint)
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

// Append stores key-value pairs in Redis with expiry
func Append(key string, value interface{}, expiryInSeconds int) error {
	conn := Pool.Get()

	conn.Send("MULTI")
	conn.Send("APPEND", key, value)

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
		if err == redigo.ErrNil {
			return reply, nil
		}

		return reply, err
	}

	return reply, err
}

// Delete deletes a value by key
func Delete(key string) error {
	conn := Pool.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}
