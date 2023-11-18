package web

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUserEmailPattern(t *testing.T) {
	// 组织测试的策略 Table Driven（表格驱动）
	testCases := []struct {
		name  string
		email string
		// 预期输出
		match bool
		//mock
		//before 数据准备
		//after 数据清洗
	}{
		{
			name:  "1",
			email: "123@",
			match: false,
		},
		{
			name:  "2",
			email: "123@qq",
			match: false,
		},
		{
			name:  "通过案例",
			email: "123@qq.com",
			match: true,
		},
	}

	handler := NewUserHandler(nil, nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := handler.emailRegexExp.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}
func TestHTTP(t *testing.T) {
	_, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("xxxx")))
	assert.NoError(t, err)

}
