package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateByID(ctx context.Context, art Article) error
}

type GormArticleDAO struct {
	db *gorm.DB
}

func (g *GormArticleDAO) UpdateByID(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := g.db.WithContext(ctx).Model(&Article{}).Where("id = ? and author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"utime":   now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
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
	Status   int
	Ctime    int64
	Utime    int64
}
