package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"webook/internal/domain"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, artId int64) (domain.Article, error)
	Set(ctx context.Context, art domain.Article) error
	GetPub(ctx context.Context, artId int64) (domain.Article, error)
	SetPub(ctx context.Context, art domain.Article) error
}

type ArticleRedisCache struct {
	cmd redis.Cmdable
}

func (r *ArticleRedisCache) GetPub(ctx context.Context, artId int64) (domain.Article, error) {
	val, err := r.cmd.Get(ctx, r.pubKey(artId)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(val, &art)
	if err != nil {
		return domain.Article{}, err
	}
	return art, nil
}

func (r *ArticleRedisCache) SetPub(ctx context.Context, art domain.Article) error {
	err := r.cmd.Set(ctx, r.pubKey(art.Id), art, time.Minute*10).Err()
	return err
}

func (r *ArticleRedisCache) Get(ctx context.Context, artId int64) (domain.Article, error) {
	val, err := r.cmd.Get(ctx, r.key(artId)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(val, &art)
	if err != nil {
		return domain.Article{}, err
	}
	return art, nil
}

func (r *ArticleRedisCache) Set(ctx context.Context, art domain.Article) error {
	err := r.cmd.Set(ctx, r.key(art.Id), art, time.Minute*10).Err()
	return err
}

func (r *ArticleRedisCache) DelFirstPage(ctx context.Context, uid int64) error {
	key := r.firstKey(uid)
	return r.cmd.Del(ctx, key).Err()
}

func (r *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := r.firstKey(uid)
	//val, err := r.cmd.Get(ctx, key).Result()
	val, err := r.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ArticleRedisCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	key := r.firstKey(uid)
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	err = r.cmd.Set(ctx, key, val, time.Minute*10).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *ArticleRedisCache) pubKey(artId int64) string {
	return fmt.Sprintf("article:pub:detail:%d", artId)
}

func (r *ArticleRedisCache) key(artId int64) string {
	return fmt.Sprintf("article:detail:%d", artId)
}

func (r *ArticleRedisCache) firstKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}

func NewArticleRedisCache(cmd redis.Cmdable) ArticleCache {
	return &ArticleRedisCache{
		cmd: cmd,
	}
}
