package crawl

import "github.com/garyburd/redigo/redis"

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

func redisDiffAndSet(rpool *redis.Pool, key, count_key string, urls []string) ([]string, error) {
	rconn := rpool.Get()
	defer rconn.Close()
	args := redis.Args{}.Add(key).Add(count_key).AddFlat(urls)

	return redis.Strings(setDiffExists.Do(rconn, args...))
}
