package dao

import (
	"context"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, artId int64, uid int64, status uint8) error
	GetByAuthor(ctx context.Context, limit, offset int, uid int64) ([]Article, error)
}
