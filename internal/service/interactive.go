package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"webook/internal/domain"
	"webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
	GetIntrByArtId(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

// GetIntrByArtId 获取文章的互动信息,包括是否点赞,是否收藏
func (i *interactiveService) GetIntrByArtId(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	intr, err := i.repo.GetInteractive(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}

	var eg = errgroup.Group{}
	eg.Go(func() error {
		intr.Liked, err = i.repo.Liked(ctx, biz, bizId, uid)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		intr.Collected, err = i.repo.Collected(ctx, biz, bizId, uid)
		if err != nil {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		// 日志
		fmt.Println(err)
	}
	return intr, nil
}

// AddCollectionItem 新增收集项
func (i *interactiveService) AddCollectionItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	return i.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.DecrLike(ctx, biz, bizId, uid)
}

func (i *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, bizId, uid)
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}
