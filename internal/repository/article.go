package repository

import (
	"context"
	"gorm.io/gorm"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDAO

	// repository层 V2分发 SyncV1专用
	authorDAO dao.ArticleAuthorDAO
	readerDAO dao.ArticleReaderDAO

	// SyncV2
	db *gorm.DB
}

func NewCacheArticleRepositoryV2(authorDAO dao.ArticleAuthorDAO, readerDAO dao.ArticleReaderDAO) *CacheArticleRepository {
	return &CacheArticleRepository{
		authorDAO: authorDAO,
		readerDAO: readerDAO,
	}
}

func NewCacheArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))
}

// SyncV2 事务开启
func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer func() {
		tx.Rollback()
	}()
	authorDAO := dao.NewGormArticleAuthorDAO(tx)
	readerDAO := dao.NewGormArticleReaderDAO(tx)
	artn := c.toEntity(art)
	var artId int64
	var err error
	if artn.Id > 0 {
		err = authorDAO.Update(ctx, artn)
	} else {
		artId, err = authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = artId
	err = readerDAO.Upsert(ctx, artn)
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return artId, nil
}

func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artn := c.toEntity(art)
	var artId int64
	var err error
	if artn.Id > 0 {
		err = c.authorDAO.Update(ctx, artn)
	} else {
		artId, err = c.authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = artId
	artPubn := dao.ArticlePublish(artn)
	err = c.readerDAO.UpsertV2(ctx, artPubn)
	if err != nil {
		return 0, err
	}
	return artId, nil
}

func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	// 更新缓存
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(art))
}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
