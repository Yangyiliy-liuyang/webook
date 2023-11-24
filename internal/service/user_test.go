package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"webook/internal/domain"
	"webook/internal/repository"
)

func TestPasswordBcrypt(t *testing.T) {
	password := []byte("121js@#hddd")
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
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.User
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &userService{
				repo: tt.fields.repo,
			}
			got, err := svc.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if !tt.wantErr(t, err, fmt.Sprintf("Login(%v, %v, %v)", tt.args.ctx, tt.args.email, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Login(%v, %v, %v)", tt.args.ctx, tt.args.email, tt.args.password)
		})
	}
}
