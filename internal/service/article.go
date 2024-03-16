package service

import (
	"context"
	"github.com/pkg/errors"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository

	// V1 专用
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository
	l          logger.Logger
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var err error
	artId := art.Id
	// 制作库
	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		artId, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = artId
	// 重置3次
	for i := 0; i < 3; i++ {
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			a.l.Error("制作库成功后，线上库保存失败", logger.Int64("artId", artId), logger.Error(err))
			return 0, err
		}
		if err == nil {
			return artId, nil
		}
	}
	a.l.Error("制作库成功后，线上库保存失败，重置次数耗尽", logger.Int64("artId", artId), logger.Error(err))
	return artId, errors.New("制作库成功后，线上库保存失败，重置次数耗尽")

}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	// 同步
	return a.repo.Sync(ctx, article)
}

func NewArticleServiceV1(authorRepo repository.ArticleAuthorRepository, readerRepo repository.ArticleReaderRepository, l logger.Logger) *articleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		l:          l,
	}
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// 借助帖子id，判断是新增还是更新
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	} else {
		return a.repo.Create(ctx, art)
	}
}
