package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/integration/startup"
	"webook/internal/repository/dao"
	ijwt "webook/internal/web/jwt"
)

// 测试套件 ArticleHandlerSuite
type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

// 设置测试前的准备
func (s *ArticleHandlerSuite) SetupSuite() {
	s.db = startup.InitDB()
	hdl := startup.InitArticleHandler()
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRouter(server)
	s.server = server
}

// 设置测试后的准备
func (s *ArticleHandlerSuite) TearDownSuite() {
	s.db.Exec("truncate table `articles`")
}

func (s *ArticleHandlerSuite) TestArticleHandler_Edit() {
	t := s.T()
	testCases := []struct {
		name     string
		befer    func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantResp Result[int64]
	}{
		{
			name:  "新建帖子",
			befer: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 验证保存到了数据库中
				var art dao.Article
				err := s.db.Where("author_id = ?", 123).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Id > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.Equal(t, "我的第一条帖子", art.Title)
				assert.Equal(t, "内容...............", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
			},
			art: Article{
				Title:   "我的第一条帖子",
				Content: "内容...............",
			},
			wantCode: 200,
			wantResp: Result[int64]{
				Data: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.befer(t)
			defer tc.after(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var resp Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&resp)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

func TestArticleHandler_EditV1(t *testing.T) {
	db := startup.InitDB()
	server := startup.InitWebServer()
	testCases := []struct {
		name     string
		befer    func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantResp Result[int64]
	}{
		{
			name:  "新建帖子",
			befer: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 验证保存到了数据库中
				var art dao.Article
				err := db.Where("author_id = ?", 123).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Id > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.Equal(t, "我的第一条帖子", art.Title)
				assert.Equal(t, "内容...............", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)

				// 清理测试数据
				db.Exec("truncate articles where author_id = ?", 123)
			},
			art: Article{
				Title:   "我的第一条帖子",
				Content: "内容...............",
			},
			wantCode: 200,
			wantResp: Result[int64]{
				Data: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.befer(t)
			defer tc.after(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var resp Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&resp)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

type Result[T any] struct {
	Success   bool   `json:"success"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Data      T      `json:"data"`
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
