package service

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
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
