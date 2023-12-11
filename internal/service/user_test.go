package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository"
	repomocks "webook/internal/repository/mocks"
)

func TestPasswordBcrypt(t *testing.T) {
	password := []byte("hello#world123455")
	//加密 bcrypt限制密码长度不超过72字节
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	//断言这个地方不应该出现err
	assert.NoError(t, err)
	println("加密后：", string(hash))
	//解密
	fakePassword := "sabdldadsaf"
	err = bcrypt.CompareHashAndPassword(hash, []byte(fakePassword))
	//有err，不为nil
	assert.NotNil(t, err)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	assert.NoError(t, err)
}

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$04$Aj5Cahd9nNle5UU/QzTNgem5bppPQcfCMe5.UyMDGGOpI2qEfqOSS",
					Phone:    "1221321312212",
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123455",
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$04$Aj5Cahd9nNle5UU/QzTNgem5bppPQcfCMe5.UyMDGGOpI2qEfqOSS",
				Phone:    "1221321312212",
			},
			wantErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123455",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123455",
			wantUser: domain.User{},
			wantErr:  errors.New("db 错误"),
		},
		{
			name: "用户名或者密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$04$Aj5Cahd9nNle5UU/QzTNgem5bppPQcfCMe5.UyMDGGOpI2qEfqOSS",
					Phone:    "1221321312212",
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world12345599999999999",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
