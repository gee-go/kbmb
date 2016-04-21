package crawl

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gee-go/kbmb/cfg"
)

type VisitCache interface {
	// DiffAndSet returns a list of urls that have not already been visited.
	// Each url is then marked as set.
	DiffAndSet(urls []string) ([][]byte, error)
	Clear() error
}

var setDiffExists = redis.NewScript(1, `
local r = {}

for _, m in pairs(ARGV) do
  if redis.call('SISMEMBER', KEYS[1], m) == 0 then
    r[#r+1] = m 
  end
end

for _, m in pairs(ARGV) do
  redis.call('SADD', KEYS[1], m)
end

return r`)

type RedisVisitCache struct {
	rpool *redis.Pool
	key   string
}

func NewRedisVisitCache() VisitCache {
	return &RedisVisitCache{
		rpool: cfg.NewRedisPool(),
		key:   "visited",
	}
}

func (vc *RedisVisitCache) DiffAndSet(urls []string) ([][]byte, error) {
	rconn := vc.rpool.Get()
	defer rconn.Close()

	return redis.ByteSlices(setDiffExists.Do(rconn, redis.Args{}.Add(vc.key).AddFlat(urls)...))
}

func (vc *RedisVisitCache) Clear() error {
	rconn := vc.rpool.Get()
	defer rconn.Close()
	_, err := rconn.Do("DEL", vc.key)
	return err
}
