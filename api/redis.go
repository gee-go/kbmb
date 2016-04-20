package api

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const redisEvalSetDiff = `
local r = {}

for _, m in pairs(ARGV) do
  if redis.call('SISMEMBER', KEYS[1], m) == 0 then
    r[#r+1] = m 
  end
end

for _, m in pairs(ARGV) do
  redis.call('SADD', KEYS[1], m)
end

return r
 `

var setDiffExists = redis.NewScript(1, redisEvalSetDiff)

// SetDiffExists returns the values that don't already exist in the set.
// Also sets the keys.
func SetDiffExists(rpool *redis.Pool, key string, values []string) ([][]byte, error) {
	rconn := rpool.Get()
	defer rconn.Close()

	return redis.ByteSlices(setDiffExists.Do(rconn, redis.Args{}.Add(key).AddFlat(values)...))
}

// NewRedisPool creates a redis pool
func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://127.0.0.1:6379")
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
