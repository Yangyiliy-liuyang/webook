package web

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	"webook/internal/service"
	svcmocks "webook/internal/service/mocks"
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
func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		//函数式编程
		//构造请求，预期中的输入
		reqBuilder func(t *testing.T) *http.Request
		//预期输出
		wantCode int
		wantBody string
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			// 利用mock构造UserHandler
			hdl := NewUserHandler(userSvc, codeSvc)
			// 注册路由
			server := gin.Default()
			hdl.RegisterRouter(server)
			// 准备HTTP请求
			req := tc.reqBuilder(t)
			// Recorder响应
			recorder := httptest.NewRecorder()
			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

//func TestHTTP(t *testing.T) {
//	// 构建http请求
//	_, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("我的请求体")))
//	assert.NoError(t, err) // 断言一定有err
//	// 获得http响应
//	recorder := httptest.NewRecorder()
//	assert.Equal(t, http.StatusOK, recorder.Code)
//}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// mock模拟实现
	userSvc := svcmocks.NewMockUserService(ctrl)
	// 设置了模拟场景，预期第一个参数是任意，第二个一定是
	userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
		Id:    1,
		Email: "1223@qq.com",
	}).Return(errors.New("db 出错"))

	err := userSvc.SingUp(context.Background(), domain.User{})
	//err := userSvc.SingUp(context.Background(), domain.User{Id: 1,
	//	Email: "1223@qq.com"})
	t.Log(err)
}
