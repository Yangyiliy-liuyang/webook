package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"webook/internal/domain"
)

const ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

type UserCacheV1 struct {
	client *redis.Client
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.Key(uid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}
func (c *RedisUserCache) Key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.Key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

// MemoryUserCache 基于本地缓存实现
type MemoryUserCache struct {
}
