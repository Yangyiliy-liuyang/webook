package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	"webook/internal/domain/proctocol"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.ArticleService
		reqBody string

		wantCode int
		wantResp proctocol.RespGeneral
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc = svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody:  `{"title":"标题","content":"内容"}`,
			wantCode: http.StatusOK,
			wantResp: proctocol.RespGeneral{
				Success:   true,
				Data:      float64(1),
				ErrorCode: 200,
				ErrorMsg:  "ok",
			},
		},
		{
			name: "已有帖子发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc = svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody:  `{"id":1,"title":"标题","content":"内容"}`,
			wantCode: http.StatusOK,
			wantResp: proctocol.RespGeneral{
				Success:   true,
				Data:      float64(1),
				ErrorCode: 200,
				ErrorMsg:  "ok",
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc = svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "标题",
					Content: "内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("发表失败"))
				return svc
			},
			reqBody:  `{"title":"标题","content":"内容"}`,
			wantCode: http.StatusOK,
			wantResp: proctocol.RespGeneral{
				Success:   true,
				ErrorCode: 500,
				ErrorMsg:  "系统内部错误",
			},
		},
		{
			name: "Bind错误",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc = svcmocks.NewMockArticleService(ctrl)
				return svc
			},
			reqBody:  `{"title":"标题","content":"内容"uuuuuuuu}`,
			wantCode: 400,
			wantResp: proctocol.RespGeneral{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			artSvc := tc.mock(ctrl)
			hdl := NewArticleHandler(artSvc, logger.NewNopLogger())
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", ijwt.UserClaims{
					Uid: 123,
				})
			})
			hdl.RegisterRouter(server)
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBufferString(tc.reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			var res proctocol.RespGeneral
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantResp, res)
		})
	}
}
