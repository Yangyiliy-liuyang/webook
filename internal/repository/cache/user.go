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

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

type UserCacheV1 struct {
	client *redis.Client
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
func (c *UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.Key(uid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}
func (c *UserCache) Key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c *UserCache) Set(ctx context.Context, du domain.User) error {
	key := c.Key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}
