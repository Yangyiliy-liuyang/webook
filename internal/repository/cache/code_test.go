package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	redismocks "webook/internal/repository/cache/rediscache"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(0), nil)
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:   context.Background(),
			phone: "15212345678",
			code:  "123456",
		},
		{
			name: "redis返回error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(0), errors.New("redis error"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: errors.New("redis error"),
		},
		{
			name: "redis返回-2，验证码不存在过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(-2), errors.New("验证码存在，但是没有过期时间"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: errors.New("验证码存在，但是没有过期时间"),
		},
		{
			name: "redis返回-1，发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(-1), errors.New("发送太频繁"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisCodeCache_Verify(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantOK  bool
		wantErr error
	}{
		{
			name: "验证通过",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(0), nil)
				cmd.EXPECT().Eval(gomock.Any(), luaVarifyCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantOK:  true,
			wantErr: nil,
		},
		{
			name: "redis返回error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(0), errors.New("redis error"))
				cmd.EXPECT().Eval(gomock.Any(), luaVarifyCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: errors.New("redis error"),
		},
		{
			name: "redis返回-2，验证码输入错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(-2), errors.New("验证码存在，验证码输入错误"))
				cmd.EXPECT().Eval(gomock.Any(), luaVarifyCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: errors.New("验证码存在，验证码输入错误"),
		},
		{
			name: "redis返回-1，验证太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewCmdResult(int64(-1), errors.New("验证太频繁"))
				cmd.EXPECT().Eval(gomock.Any(), luaVarifyCode, gomock.Any(), gomock.Any()).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			phone:   "15212345678",
			code:    "123456",
			wantErr: ErrCodeVerifyTooMany,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCodeCache(tc.mock(ctrl))
			ok, err := c.Verify(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
