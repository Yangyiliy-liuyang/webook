package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

var (
	ErrDuplicateUser = repository.ErrDuplicateUser
	// todo 含糊
	ErrInvalidUserOrPassword = errors.New("用户或者密码不正确")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (svc *UserService) SingUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *UserService) UpdateUserInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateUserInfo(ctx, user)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 兼顾性能
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		//err!=nil 系统错误 或者 err==nil 已经找到
		return u, err
	}
	//没找到，注册新用户
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil || !errors.Is(err, ErrDuplicateUser) {
		// 有错误，但是不是phone唯一索引错误   系统错误
		return domain.User{}, err
	}
	// 用户存在 err == nil || errors.Is(err, ErrDuplicateUser)
	// todo 主从延迟 不一定能在数据库中找到数据
	// 插入的主库，查询的是从库
	// 强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserService) GetUserInfo(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
