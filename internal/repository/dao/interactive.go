package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

// InsertLikeInfo
func (g *GormInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Upsert
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			BizId:  bizId,
			Biz:    biz,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}
		// Upsert likeCnt
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("like_cnt + 1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			BizId:   bizId,
			Biz:     biz,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (g *GormInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 必然存在，所以Update 即可
		err := tx.Model(UserLikeBiz{}).
			Where("uid = ? and biz_id = ? and biz = ? ", uid, bizId, biz).
			Updates(map[string]interface{}{
				"status": 0,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}
		// Update likeCnt
		return tx.Model(&Interactive{}).Where("biz_id = ? and biz = ? ", bizId, biz).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("like_cnt - 1"),
				"utime":    now,
			}).Error
	})
}

func (g *GormInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	// 新帖子，没有统计，插入数据
	// 更新帖子，更新阅读数
	// Upsert语义

	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			// 更新数据不能先取出原数据加一，再更新回数据库，会导致并发问题
			"read_cnt": gorm.Expr("read_cnt + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		BizId:   bizId,
		Biz:     biz,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

func NewGormInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GormInteractiveDAO{
		db: db,
	}
}

type Interactive struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 唯一索引 <bizId,biz>
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"uniqueIndex:biz_type_id"`
	ReadCnt    int64  // 阅读数
	LikeCnt    int64  //点赞数
	CollectCnt int64  //收藏数
	CommentCnt int64  //评论数
	ShareCnt   int64  // 分享数
	Ctime      int64
	Utime      int64
}

type UserLikeBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 唯一索引 <bizId,biz>
	Uid    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"uniqueIndex:uid_biz_type_id"`
	Status uint
	Ctime  int64
	Utime  int64
}
