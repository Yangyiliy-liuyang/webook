package repository

import (
	"context"
	"webook/internal/domain"
)

type ArticleReaderRepository interface {
	// Save 相当于insert or create
	Save(ctx context.Context, art domain.Article) error
}
