package dao

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
	"webook/internal/domain"
)

var statusPrivate = domain.ArticleStatusPrivate.ToUint8()

type S3DAO struct {
	oss    *s3.S3
	bucket *string
	GormArticleDAO
}

func (d *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	var artId int64
	var err error
	err = d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		pubArt := ArticlePublish(art)
		err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
			DoNothing: false,
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  pubArt.Title,
				"utime":  now,
				"status": pubArt.Status,
			}),
			UpdateAll: false,
		}).Create(&pubArt).Error
		return err
	})
	d.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      d.bucket,
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return artId, err
}

func (d *S3DAO) SyncStatus(ctx context.Context, artId int64, uid int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		return err
	}

	if status == statusPrivate {
		_, err = d.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: d.bucket,
			Key:    ekit.ToPtr[string](strconv.FormatInt(uid, 10)),
		})
	}
	return err
}

func NewOssDAO(oss *s3.S3, db *gorm.DB) ArticleDAO {
	return &S3DAO{
		oss:    oss,
		bucket: ekit.ToPtr[string]("webook-1314583317"),
		GormArticleDAO: GormArticleDAO{
			db: db,
		},
	}
}
