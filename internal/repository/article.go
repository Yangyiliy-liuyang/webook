package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

func NewCacheArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.ID,
	}
}
