# redis
--
    import "github.com/reddotpay/redis"


## Usage

```go
var Pool *redigo.Pool
```
Pool contains Redis pool

#### func  Delete

```go
func Delete(key string) error
```
Delete deletes a value by key

#### func  NewPool

```go
func NewPool(endpoint string, maxIdle int, idleTimeout time.Duration) *redigo.Pool
```
NewPool initialises a new Redis pool

#### func  Retrieve

```go
func Retrieve(key string) (interface{}, error)
```
Retrieve retrieves value by key

#### func  Store

```go
func Store(key string, value interface{}, expiryInSeconds int) error
```
Store stores key-value pairs in Redis with expiry
