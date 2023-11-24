package web

import (
	"bytes"
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

		// mock
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 构造请求，预期中输入
		reqBuilder func(t *testing.T) *http.Request

		// 预期中的输出
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "hello,正在注册...",
		},
		{
			name: "Bind出错",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusBadRequest,
		},

		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "密码格式错误，必须包含字母、数字、特殊字符",
		},
		{
			name: "两次密码输入不同",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123455",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "两次输入密码不对",
		},

		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello",
"confirmPassword": "hello"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "密码必须包含字母、数字、特殊字符，并且不少于八位",
		},

		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("db错误"))
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},

		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrDuplicateUser)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "邮箱冲突，请换一个",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造 handler
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			// 准备服务器，注册路由
			server := gin.Default()
			hdl.RegisterRouter(server)

			// 准备Req和记录的 recorder
			req := tc.reqBuilder(t)
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
//
//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	// mock模拟实现
//	userSvc := svcmocks.NewMockUserService(ctrl)
//	// 设置了模拟场景，预期第一个参数是任意，第二个一定是
//	userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
//		Id:    1,
//		Email: "1223@qq.com",
//	}).Return(errors.New("db 出错"))
//
//	err := userSvc.SingUp(context.Background(), domain.User{})
//	//err := userSvc.SingUp(context.Background(), domain.User{Id: 1,
//	//	Email: "1223@qq.com"})
//	t.Log(err)
//}
