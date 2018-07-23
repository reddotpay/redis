package redis_test

import (
	"errors"
	"testing"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/reddotpay/redis"
	"github.com/stretchr/testify/assert"
)

var (
	mockRedisPool *redigo.Pool
	mockRedisConn *redigomock.Conn
)

func init() {
	mockRedisConn = redigomock.NewConn()
	redis.Pool = redigo.NewPool(func() (redigo.Conn, error) {
		return mockRedisConn, nil
	}, 10)
}

func TestRedis_NewPool(t *testing.T) {
	pool := redis.NewPool("endpoint", 1, 1)
	defer pool.Close()
	assert.IsType(t, redis.Pool, pool)
}

func TestRedis_Store_Bytes(t *testing.T) {
	expiryInSeconds := 10
	mockRedisConn.Clear()
	mockRedisConn.Command("MULTI")
	mockRedisConn.Command("SET", "key", []byte("value"))
	mockRedisConn.Command("EXPIRE", "key", expiryInSeconds)
	mockRedisConn.Command("EXEC")
	err := redis.Store("key", []byte("value"), expiryInSeconds)
	assert.Nil(t, err)
}

func TestRedis_Store_String(t *testing.T) {
	expiryInSeconds := 10
	mockRedisConn.Clear()
	mockRedisConn.Command("MULTI")
	mockRedisConn.Command("SET", "key", "value")
	mockRedisConn.Command("EXPIRE", "key", expiryInSeconds)
	mockRedisConn.Command("EXEC")
	err := redis.Store("key", "value", expiryInSeconds)
	assert.Nil(t, err)
}

func TestRedis_Store_NoExpiry(t *testing.T) {
	expiryInSeconds := 0
	mockRedisConn.Clear()
	mockRedisConn.Command("MULTI")
	mockRedisConn.Command("SET", "key", "value")
	mockRedisConn.Command("EXEC")
	err := redis.Store("key", "value", expiryInSeconds)
	assert.Nil(t, err)
}

func TestRedis_Store_Error(t *testing.T) {
	expiryInSeconds := 0
	mockRedisConn.Clear()
	mockRedisConn.Command("MULTI")
	mockRedisConn.Command("SET", "key", "value").ExpectError(errors.New("ERR"))
	mockRedisConn.Command("EXEC")
	err := redis.Store("key", "value", expiryInSeconds)
	assert.NotNil(t, err)
	assert.Equal(t, "ERR", err.Error())
}

func TestRedis_Retrieve_Bytes(t *testing.T) {
	mockRedisConn.Clear()
	mockRedisConn.Command("GET", "key").Expect([]byte("value"))
	reply, err := redis.Retrieve("key")
	assert.Nil(t, err)
	assert.Equal(t, []byte("value"), reply.([]byte))
}

func TestRedis_Retrieve_String(t *testing.T) {
	mockRedisConn.Clear()
	mockRedisConn.Command("GET", "key").Expect("value")
	reply, err := redis.Retrieve("key")
	assert.Nil(t, err)
	assert.Equal(t, "value", reply.(string))
}

func TestRedis_Retrieve_NotFound(t *testing.T) {
	mockRedisConn.Clear()
	mockRedisConn.Command("GET", "key").ExpectError(redigo.ErrNil)
	reply, err := redis.Retrieve("key")
	assert.Nil(t, err)
	assert.Nil(t, reply)
}

func TestRedis_Retrieve_Error(t *testing.T) {
	mockRedisConn.Clear()
	mockRedisConn.Command("GET", "key").ExpectError(errors.New("ERR"))
	reply, err := redis.Retrieve("key")
	assert.NotNil(t, err)
	assert.Equal(t, "ERR", err.Error())
	assert.Empty(t, reply)
}
