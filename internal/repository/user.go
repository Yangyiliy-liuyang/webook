package repository

import (
	"context"
	"database/sql"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	// todo repository是与业务强相关的，是一定要
	ErrUserNotFound = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, OpenId string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CacheUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	du, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	u := repo.toDomain(du)
	return u, nil
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	du, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	u := repo.toDomain(du)
	return u, nil
}

func (repo *CacheUserRepository) toDomain(du dao.User) domain.User {
	return domain.User{
		Id:       du.Id,
		Email:    du.Email.String,
		Password: du.Password,
		Phone:    du.Phone.String,
		AboutMe:  du.AboutMe,
		Nickname: du.Nickname,
		Birthday: time.UnixMilli(du.Birthday),
		Ctime:    time.UnixMilli(du.Ctime),
		WechatInfo: domain.WechatInfo{
			OpenId:  du.WetchatOpenId.String,
			UnionId: du.WetchatUnionId.String,
		},
	}
}

func (repo *CacheUserRepository) UpdateUserInfo(ctx context.Context, u domain.User) error {
	return repo.dao.InsertInfo(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			// false 为空 true 为 不为空0
			Valid: u.Phone != "",
		},
		Birthday:       u.Birthday.UnixMilli(),
		Nickname:       u.Nickname,
		AboutMe:        u.AboutMe,
		Password:       u.Password,
		WetchatOpenId:  sql.NullString{String: u.WechatInfo.OpenId, Valid: u.WechatInfo.OpenId != ""},
		WetchatUnionId: sql.NullString{String: u.WechatInfo.UnionId, Valid: u.WechatInfo.UnionId != ""},
	}
}

func (repo *CacheUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, uid)
	// 注意这里的处理方式
	if err == nil {
		return u, err
	}
	du, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	u = repo.toDomain(du)
	// 忽略掉这里的错误
	_ = repo.cache.Set(ctx, u)
	return u, nil
}

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	du, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	u := repo.toDomain(du)
	return u, nil
}
