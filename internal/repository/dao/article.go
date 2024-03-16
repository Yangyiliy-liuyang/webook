package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
}

type GormArticleDAO struct {
	db *gorm.DB
}

func (g *GormArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	tx := g.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() {
		tx.Rollback()
	}()
	var artId int64
	var err error
	dao := NewGormArticleDAO(tx)
	if art.Id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		artId, err = dao.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = artId

	pubArt := ArticlePublish(art)
	err = tx.Clauses(clause.OnConflict{
		DoNothing: false,
		DoUpdates: clause.Assignments(map[string]interface{}{
,			"title":   pubArt.Title,
            "content": pubArt.Content,
            "utime":   pubArt.Utime,
		}),
		UpdateAll: false,
	}).Create(&pubArt).Error

	if err != nil {
		return 0, err
	}
	tx.Commit()
	return artId, nil
}

func (g *GormArticleDAO) UpdateById(ctx context.Context, art Article) error {
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

// ArticlePublish 线上库表
type ArticlePublish Article

// ArticlePublish2 写法二
type ArticlePublish2 struct {
	Article
}
