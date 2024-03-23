package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"webook/internal/domain"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const filedReadCnt = "read_cnt"
const filedLikeCnt = "like_cnt"
const filedCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	GetInteractive(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	SetInteractive(ctx context.Context, res domain.Interactive) error
}

type InteractiveRedisCache struct {
	cmd redis.Cmdable
}

func (i *InteractiveRedisCache) SetInteractive(ctx context.Context, res domain.Interactive) error {
	key := i.key(res.Biz, res.BizId)
	err := i.cmd.HMSet(ctx, key, map[string]interface{}{
		filedReadCnt:    res.ReadCnt,
		filedLikeCnt:    res.LikeCnt,
		filedCollectCnt: res.CollectCnt,
	}).Err()
	if err != nil {
		return err
	}

	// 设置过期时间
	return i.cmd.Expire(ctx, key, 15*time.Minute).Err()
}

func (i *InteractiveRedisCache) GetInteractive(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	key := i.key(biz, bizId)
	res, err := i.cmd.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	readCnt, _ := strconv.ParseInt(res[filedReadCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(res[filedLikeCnt], 10, 64)
	collectCnt, _ := strconv.ParseInt(res[filedCollectCnt], 10, 64)
	return domain.Interactive{
		Biz:        biz,
		BizId:      bizId,
		ReadCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
	}, nil
}

func (i *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.cmd.Eval(ctx, luaIncrCnt, []string{key}, filedCollectCnt, 1).Err()
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
