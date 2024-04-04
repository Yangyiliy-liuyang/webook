package service

import (
	"context"
	"webook/internal/domain"
)

type RankingService interface {
	TopN(ctx context.Context) error // 前一百
}

type BatchRankingService struct {
}

func (b BatchRankingService) TopN(ctx context.Context) error {
	arts, err := b.topN(ctx)
	if err != nil {
		return err
	}
	//TODO implement me
	panic(arts)
	// 最终是放在缓存里面的
}

func (b BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {

}
