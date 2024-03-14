package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type GormArticleDAO struct {
	db *gorm.DB
}

func NewGormArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{
		db: db,
	}
}

func (g *GormArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

type Article struct {
	Id      int64  `gorm:"primary_key autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type:BLOB"`
	// 索引
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}
