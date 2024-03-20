package repository

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, artId int64, uid int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, limit, offset int, uid int64) ([]domain.Article, error)
	GetByArtId(ctx context.Context, artId int64) (domain.Article, error)
}

func (c CacheArticleRepository) GetByArtId(ctx context.Context, artId int64) (domain.Article, error) {
	res, err := c.cache.Get(ctx, artId)
	if err == nil {
		return res, err
	}
	art, err := c.dao.GetByArtId(ctx, artId)
	if err != nil {
		return domain.Article{}, err
	}
	go func() {
		err := c.cache.Set(ctx, c.toDomain(art))
		if err != nil {
			// 记录日志 监控
			fmt.Println("缓存回写失败", err)
		}
	}()
	return c.toDomain(art), nil
}

func (c *CacheArticleRepository) GetByAuthor(ctx context.Context, limit, offset int, uid int64) ([]domain.Article, error) {
	if limit == 100 && offset == 0 {
		res, err := c.cache.GetFirstPage(ctx, uid)
		if res == nil {
			return res, nil
		} else {
			// 缓存未命中，记录日志
			return res, err
		}
	}
	arts, err := c.dao.GetByAuthor(ctx, limit, offset, uid)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		// 缓存回写失败，有可能是大问题 网络问题连不上redis 或者第三方redis问题
		err := c.cache.SetFirstPage(ctx, uid, res)
		if err != nil {
			// 记录日志 监控
			fmt.Println("缓存回写失败", err)
		}
	}()
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		c.preCache(ctx, res)
	}()
	return res, nil
}

type CacheArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
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

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, artId int64, uid int64, status domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, artId, uid, status.ToUint8())
	if err != nil {
		return err
	}
	err = c.cache.DelFirstPage(ctx, uid)
	if err != nil {
		return err
	}
	return nil
}

func NewCacheArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &CacheArticleRepository{
		dao:   dao,
		cache: cache,
	}
}

func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	artId, err := c.dao.Sync(ctx, c.toEntity(art))
	if err != nil {
		return 0, err
	}
	err = c.cache.DelFirstPage(ctx, artId)
	if err != nil {
		return 0, err
	}
	return artId, nil

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
	err := c.dao.UpdateById(ctx, c.toEntity(art))
	if err != nil {
		return err
	}
	err = c.cache.DelFirstPage(ctx, art.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	artId, err := c.dao.Insert(ctx, c.toEntity(art))
	if err != nil {
		return 0, err
	}
	err = c.cache.DelFirstPage(ctx, artId)
	if err != nil {
		return 0, err
	}
	return artId, nil
}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		//Status:   uint8(art.Status),
		// 连调写法
		Status: art.Status.ToUint8(),
	}
}

func (c *CacheArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id:   art.AuthorId,
			Name: art.AuthorName,
		},
		Status: domain.ArticleStatus(art.Status),
		Ctime:  art.Ctime,
		Utime:  art.Utime,
	}
}

func (c *CacheArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	// 缓存回写
	const size = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) < size {
		err := c.cache.Set(ctx, arts[0])
		if err != nil {
			// 记录日志 监控
			fmt.Println("缓存回写失败", err)
		}
	}
}
