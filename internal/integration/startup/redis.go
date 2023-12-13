package startup

import "github.com/go-redis/redis/v8"

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
