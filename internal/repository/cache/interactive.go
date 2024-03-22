package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const filedReadCnt = "read_cnt"
const filedLikeCnt = "like_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
}

type InteractiveRedisCache struct {
	cmd redis.Cmdable
}

func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.cmd.Eval(ctx, luaIncrCnt, []string{key}, filedLikeCnt, 1).Err()
}

func (i *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.cmd.Eval(ctx, luaIncrCnt, []string{key}, filedLikeCnt, -1).Err()
}

func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	// 不需要 0 的返回值
	return i.cmd.Eval(ctx, luaIncrCnt, []string{key}, filedReadCnt, 1).Err()
}

func (i *InteractiveRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:article:%s:%d", biz, bizId)
}

func NewInteractiveCache(cmd redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		cmd: cmd,
	}
}
