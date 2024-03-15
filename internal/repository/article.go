package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	UpdateByID(ctx context.Context, art domain.Article) error
}

type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

func NewCacheArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

func (c *CacheArticleRepository) UpdateByID(ctx context.Context, art domain.Article) error {
	// 更新缓存
	return c.dao.UpdateByID(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
