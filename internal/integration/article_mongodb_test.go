package integration

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/internal/integration/startup"
	"webook/internal/repository/dao"
	ijwt "webook/internal/web/jwt"
)

// 测试套件 ArticleMongoHandlerSuite
type ArticleMongoHandlerSuite struct {
	suite.Suite
	mdb     *mongo.Database
	col     *mongo.Collection // 制作库
	liveCol *mongo.Collection //线上库
	server  *gin.Engine
}

// 设置测试前的准备
func (s *ArticleMongoHandlerSuite) SetupSuite() {
	s.mdb = startup.InitMongoDB()
	s.col = s.mdb.Collection("articles")
	s.liveCol = s.mdb.Collection("published_articles")
	node, err := snowflake.NewNode(1)
	assert.NoError(s.T(), err)
	hdl := startup.InitArticleHandler(dao.NewMongoDBArticleDAO(s.mdb, node))
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
func (s *ArticleMongoHandlerSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.col.DeleteMany(ctx, bson.M{})
	assert.NoError(s.T(), err)
	_, err = s.liveCol.DeleteMany(ctx, bson.M{})
	assert.NoError(s.T(), err)
}

func (s *ArticleMongoHandlerSuite) TestArticleHandler_Edit() {
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
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				// 验证保存到了数据库中
				var art dao.Article
				err := s.col.FindOne(ctx, bson.D{bson.E{"author_id", 123}}).Decode(&art)
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
				Success:   true,
				Data:      1,
				ErrorCode: 200,
				ErrorMsg:  "ok",
			},
		},
		//{
		//	name: "修改帖子",
		//	befer: func(t *testing.T) {
		//		err := s.db.Create(&dao.Article{
		//			Id:       2,
		//			Title:    "我的帖子2",
		//			Content:  "内容........2......",
		//			AuthorId: 123,
		//			Ctime:    163744444,
		//			Utime:    1637444444,
		//		}).Error
		//		assert.NoError(t, err)
		//	},
		//	after: func(t *testing.T) {
		//		art := dao.Article{}
		//		err := s.db.Where("id = ?", 2).First(&art).Error
		//		assert.NoError(t, err)
		//		// 验证保存到了数据库中
		//		assert.True(t, art.Utime > 163744444)
		//		art.Utime = 0
		//		assert.Equal(t, dao.Article{
		//			Id:       2,
		//			Title:    "我的帖子2...修改版",
		//			Content:  "内容........2.......",
		//			AuthorId: 123,
		//			Ctime:    163744444,
		//		}, art)
		//	},
		//	art: Article{
		//		Id:      2,
		//		Title:   "我的帖子2...修改版",
		//		Content: "内容........2.......",
		//	},
		//	wantCode: 200,
		//	wantResp: Result[int64]{
		//		Success:   true,
		//		Data:      2,
		//		ErrorCode: 200,
		//		ErrorMsg:  "ok",
		//	},
		//},
		//{
		//	name: "违法修改别人的帖子",
		//	befer: func(t *testing.T) {
		//		err := s.db.Create(&dao.Article{
		//			Id:       3,
		//			Title:    "我的帖子3",
		//			Content:  "内容........3......",
		//			AuthorId: 223,
		//			Ctime:    1637444444,
		//			Utime:    1637444444,
		//		}).Error
		//		assert.NoError(t, err)
		//	},
		//	after: func(t *testing.T) {
		//		// 验证数据修改未成功
		//		art := dao.Article{}
		//		err := s.db.Where("id = ?", 3).First(&art).Error
		//		assert.NoError(t, err)
		//		assert.Equal(t, dao.Article{
		//			Id:       3,
		//			Title:    "我的帖子3",
		//			Content:  "内容........3......",
		//			AuthorId: 223,
		//			Ctime:    1637444444,
		//			Utime:    1637444444,
		//		}, art)
		//	},
		//	art: Article{
		//		Id:      3,
		//		Title:   "我的帖子3...修改版",
		//		Content: "内容........3.......",
		//	},
		//	wantCode: 200,
		//	wantResp: Result[int64]{
		//		Success:   true,
		//		ErrorCode: 500,
		//		ErrorMsg:  "系统内部错误",
		//	},
		//},
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
			if tc.wantResp.Data >= 0 {
				assert.True(t, resp.Data > 0)
			}
		})
	}
}

func TestArticleMongoHandler(t *testing.T) {
	suite.Run(t, &ArticleMongoHandlerSuite{})
}
