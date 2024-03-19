package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GormArticleDAO struct {
	db *gorm.DB
}

// GetByAuthor 根据作者ID获取文章列表
func (g *GormArticleDAO) GetByAuthor(ctx context.Context, limit, offset int, uid int64) ([]Article, error) {
	var arts []Article
	err := g.db.WithContext(ctx).Model(&Article{}).Where("author_id = ?", uid).
		Limit(limit).Offset(offset).Order("utime desc").
		Find(&arts).Error
	if err != nil {
		return nil, err
	}
	return arts, nil
}

func (g *GormArticleDAO) SyncStatus(ctx context.Context, artId int64, uid int64, status uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? and author_id = ?", artId, uid).Updates(map[string]interface{}{
			"status": status,
			"utime":  now,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("没有修改权限，更新失败")
		}
		return tx.Model(&ArticlePublish{}).Where("id = ? ", artId).Updates(map[string]interface{}{
			"status": status,
			"utime":  now,
		}).Error
	})
}

// Sync 闭包事务
func (g *GormArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var artId int64
	var err error
	err = g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 基于事务实现dao的复用
		dao := NewGormArticleDAO(tx)
		if art.Id > 0 {
			err = dao.UpdateById(ctx, art)
		} else {
			artId, err = dao.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = artId
		now := time.Now().UnixMilli()
		art.Ctime = now
		art.Utime = now
		// 衍生类型写法
		pubArt := ArticlePublish(art)

		/*
			// 使用结构体标签 衍生类型写法2
			pubArt2 := ArticlePublish2{
				Article: art,
			}
		*/

		err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: false,
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   pubArt.Title,
				"content": pubArt.Content,
				"utime":   now,
				"status":  pubArt.Status,
			}),
			UpdateAll: false,
		}).Create(&pubArt).Error
		return err
	})
	return artId, err
}

func (g *GormArticleDAO) SyncV1(ctx context.Context, art Article) (int64, error) {
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
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	pubArt := ArticlePublish(art)
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
		DoNothing: false,
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   pubArt.Title,
			"content": pubArt.Content,
			"utime":   now,
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
		"status":  art.Status,
		"utime":   now,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("没有修改权限，更新失败")
	}
	return nil
}

func (g *GormArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func NewGormArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{
		db: db,
	}
}

type Article struct {
	Id      int64  `gorm:"primary_key autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type:BLOB" bson:"content,omitempty"`
	// 索引
	AuthorId   int64  `gorm:"index" bson:"author_id,omitempty"`
	AuthorName string `gorm:"author_name" bson:"author_name,omitempty"`
	Status     uint8  ` bson:"status,omitempty"`
	Ctime      int64  `bson:"ctime,omitempty"`
	Utime      int64  `bson:"utime,omitempty"`
}

// ArticlePublish 线上库表
type ArticlePublish Article

// ArticlePublish2 写法二
type ArticlePublish2 struct {
	Article
}
