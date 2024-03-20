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
	Withdraw(ctx context.Context, artId int64, id int64) error
	GetByAuthor(ctx context.Context, limit, offset int, uid int64) ([]domain.Article, error)
	GetByArtId(ctx context.Context, artId int64) (domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository

	// V1 专用
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository
	l          logger.Logger
}

func (a *articleService) GetByArtId(ctx context.Context, artId int64) (domain.Article, error) {
	return a.repo.GetByArtId(ctx, artId)
}

func (a *articleService) GetByAuthor(ctx context.Context, limit, offset int, uid int64) ([]domain.Article, error) {
	return a.repo.GetByAuthor(ctx, limit, offset, uid)
}

func (a *articleService) Withdraw(ctx context.Context, artId int64, id int64) error {
	return a.repo.SyncStatus(ctx, artId, id, domain.ArticleStatusPrivate)
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

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	// 同步
	return a.repo.Sync(ctx, art)
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
	art.Status = domain.ArticleStatusUnPublished
	// 借助帖子id，判断是新增还是更新
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	} else {
		return a.repo.Create(ctx, art)
	}
}
