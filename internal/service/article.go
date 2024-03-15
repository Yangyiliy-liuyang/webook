package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository

	// V1 专用
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleRepository
}

func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	return 0, nil
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	return 0, nil
}

func NewArticleServiceV1(authorRepo repository.ArticleAuthorRepository, readerRepo repository.ArticleRepository) *articleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
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
		err := a.repo.UpdateByID(ctx, art)
		return art.Id, err
	} else {
		return a.repo.Create(ctx, art)
	}
}
