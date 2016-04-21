package crawl

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gee-go/kbmb/cfg"
)

type VisitCache interface {
	// DiffAndSet returns a list of urls that have not already been visited.
	// Each url is then marked as set.
	DiffAndSet(urls []string) ([][]byte, error)
	Clear() error

	// Count of unprocessed keys. Basically a distributed wait group.
	Count() (int, error)

	// Mark a key as processed. See sync.WaitGroup
	Done() error

	// Wait blocks until all keys are processed
	Wait() error
}

var setDiffExists = redis.NewScript(2, `
local r = {}

for _, m in pairs(ARGV) do
  if redis.call('SISMEMBER', KEYS[1], m) == 0 then
    r[#r+1] = m
    redis.call('SADD', KEYS[1], m)
    redis.call('INCR', KEYS[2])
  end
end

return r`)

// RedisVisitCache is a distributed VisitCache based on redis.
type RedisVisitCache struct {
	rpool     *redis.Pool
	key       string
	count_key string
}

func NewRedisVisitCache(key string) VisitCache {
	return &RedisVisitCache{
		rpool:     cfg.NewRedisPool(),
		key:       key,
		count_key: key + "_count",
	}
}

func (vc *RedisVisitCache) DiffAndSet(urls []string) ([][]byte, error) {
	rconn := vc.rpool.Get()
	defer rconn.Close()

	return redis.ByteSlices(setDiffExists.Do(rconn, redis.Args{}.Add(vc.key).Add(vc.count_key).AddFlat(urls)...))
}

// Done, just like sync.WaitGroup done.
func (vc *RedisVisitCache) Done() error {
	rconn := vc.rpool.Get()
	defer rconn.Close()
	_, err := rconn.Do("DECR", vc.count_key)
	return err
}

func (vc *RedisVisitCache) Count() (int, error) {
	rconn := vc.rpool.Get()
	defer rconn.Close()
	return redis.Int(rconn.Do("GET", vc.count_key))
}

func (vc *RedisVisitCache) Clear() error {
	rconn := vc.rpool.Get()
	defer rconn.Close()
	_, err := rconn.Do("DEL", vc.key, vc.count_key)
	return err
}

func (vc *RedisVisitCache) Wait() error {
	for range time.Tick(100 * time.Millisecond) {
		count, err := vc.Count()
		if err != nil || count == 0 {
			return err
		}
	}

	return nil
}
