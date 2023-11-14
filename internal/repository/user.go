package repository

import (
	"context"
	"database/sql"
	"log"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	// todo repository是与业务强相关的，是一定要
	ErrUserNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	du := repo.toDomain(u)
	return du, nil
}

func (repo *UserRepository) toDomain(du dao.User) domain.User {
	return domain.User{
		Id:       du.Id,
		Email:    du.Email.String,
		Password: du.Password,
	}
}

func (repo *UserRepository) UpdateUserInfo(ctx context.Context, u domain.User) error {
	return repo.dao.InsertInfo(ctx, repo.domainToEntity(u))
}
func (repo *UserRepository) domainToEntity(u domain.User) dao.User {
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
		Birthday: u.Birthday.UnixMilli(),
		Nickname: u.Nickname,
		AboutMe:  u.AboutMe,
		Password: u.Password,
	}
}
func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	//只要err为nil就返回
	//if err == nil {
	//	// 拿到了数据
	//	return du, nil
	//}
	/*switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		u, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		//同步写
		//err = repo.cache.Set(ctx, du)
		//异步写 进一步提高性能
		go func() {
			err := repo.cache.Set(ctx, du)
			if err != nil {
				log.Println(err)
			}
		}()
	default:
		//redis有问题 缓存穿透击穿，保住数据库
		return domain.User{}, err
	}*/
	//不为nil，查询数据库
	//err存在情况
	//1、key不存在，redis正常
	//2.访问redis有问题，网络 、 redis本身有问题
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	//同步写
	//err = repo.cache.Set(ctx, du)
	//异步写 进一步提高性能
	go func() {
		err := repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}
	}()
	return du, nil

}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	du, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	u := repo.toDomain(du)
	return u, nil
}
