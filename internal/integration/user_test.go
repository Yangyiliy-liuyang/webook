package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/internal/integration/startup"
	"webook/internal/web"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}
func TestUserHandler_SendSMSCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		phone    string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:17355557222"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				duration, err := rdb.TTL(ctx, key).Result()
				assert.True(t, duration > time.Minute*9+time.Second+50)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "17355557222",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 200,
				Msg:  "发送成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewReader([]byte(fmt.Sprintf(`{"phone"= "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)

			if tc.wantCode != http.StatusOK {
				return
			}
			// 反序列化为结构体
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
