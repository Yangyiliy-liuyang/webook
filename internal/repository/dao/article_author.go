package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type GormArticleAuthorDAO struct {
	db *gorm.DB
}

func NewGormArticleAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &GormArticleAuthorDAO{
		db: db,
	}
}

func (g *GormArticleAuthorDAO) Create(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleAuthorDAO) Update(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}
