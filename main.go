package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/middleware"
	"webook/internal/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHandler(db, server)
	server.Run(":8080")

}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRouter(server)
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3030"},
		AllowHeaders:     []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			//if strings.HasPrefix(origin,"http://localhost") {
			if strings.Contains(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "公司域名.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	login := &middleware.LoginMiddlewareBuilder{}
	// todo 存储数据的 userId直接存在cookie中
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("userId", store), login.CheckLogin())
	return server
}

func initDB() *gorm.DB {
	dsn := "root:root@tcp(localhost:13316)/webook"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
