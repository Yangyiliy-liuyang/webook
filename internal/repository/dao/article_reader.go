package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	// Upsert INSERT or UPDATE
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art ArticlePublish) error
}

type GormArticleReaderDAO struct {
	db *gorm.DB
}

func (g *GormArticleReaderDAO) UpsertV2(ctx context.Context, art ArticlePublish) error {
	//TODO implement me
	panic("implement me")
}

func NewGormArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &GormArticleReaderDAO{
		db: db,
	}
}

func (g *GormArticleReaderDAO) Upsert(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}
