package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
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
		AllowOrigins: []string{"http://localhost:3030"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//允许前端访问后端响应中带的头部
		ExposeHeaders:    []string{"X-Jwt-Token"},
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
	//  单机单实例部署 考虑基于内存的memstore实现，多实例部署，redis
	//store := cookie.NewStore([]byte("secret"))
	//memstore 是基于内存实现的
	//参数一是authentication key 身份验证
	//参数二encryption key 数据加密
	// + 授权(权限控制)就是信息安全的三个核心概念
	//最好64或者32位
	//百度 -》 生成密码 -》 复制粘贴
	//store := memstore.NewStore([]byte("05kcS4LEzQcepWhpjjas07YNzJgxL93a"),
	//	[]byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgX"))
	store, err := redis.NewStore(16, "tcp", "localhost:6379",
		"", []byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgX"),
		[]byte("05kcS4LEzQcepWhpjjas07YNzJgxL93a"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
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
