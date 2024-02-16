package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//db := initDB()
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//})
	//server := initWebServer()
	//codeSvc := initCodeSvc(redisClient)
	//initUserHandler(db, redisClient, codeSvc, server)
	//server := gin.Default()
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	server.Run(":8080")
}

//func initUserHandler(db *gorm.DB, redisClient redis.Cmdable, codeSvc *service.CodeService, server *gin.Engine) {
//	ud := dao.NewUserDAO(db)
//	uc := cache.NewUserCache(redisClient)
//	ur := repository.NewUserRepository(ud, uc)
//	us := service.NewUserService(ur)
//	hdl := web.NewUserHandler(us, codeSvc)
//	hdl.RegisterRouter(server)
//}
//
//func initCodeSvc(redisClient redis.Cmdable) *service.CodeService {
//	cc := cache.NewCodeCache(redisClient)
//	cr := repository.NewCodeRepository(cc)
//	// todo
//	return service.NewCodeService(cr, initMemorySms())
//}
//
//func initMemorySms() *localsms.Service {
//	return localsms.NewService()
//}
//
//func initWebServer() *gin.Engine {
//	server := gin.Default()
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins: []string{"http://localhost:3030"},
//		AllowHeaders: []string{"Content-Type", "Authorization"},
//		//允许前端访问后端响应中带的头部
//		ExposeHeaders:    []string{"X-Jwt-Token"},
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			//if strings.HasPrefix(origin,"http://localhost") {
//			if strings.Contains(origin, "http://localhost") {
//				return true
//			}
//			return strings.Contains(origin, "公司域名.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//	// todo 限流 err
//	//redisClient := redis.NewClient(&redis.Options{
//	//	Addr: "localhost:6379",
//	//})
//	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//	//useSession(server)
//	useJWT(server)
//	return server
//}
//
//func useJWT(server *gin.Engine) {
//	login := &middleware.LoginJWTMilddlewareBuilder{}
//	server.Use(login.CheckLoginJWT())
//}
//
//func useSession(server *gin.Engine) {
//	//login := &middleware.LoginMiddlewareBuilder{}
//	/*	//通过cookie.NewStore方法创建了一个存储用户会话数据的cookie存储实例，
//		//并使用一个密钥([]byte(“dasfdadFDFDDAS”))对会话数据进行加密。
//		store := cookie.NewStore([]byte("dasfdadFDFDDAS"))
//		//使用sessions.Sessions方法将会话中间件添加到路由器（router）中，
//		//将会话ID设置为"ssid"，并将会话存储设置为之前创建的cookie存储实例。
//		server.Use(sessions.Sessions("ssid", store))
//		//使用router.Use方法将登录中间件添加到路由器中，以在每个请求处理之前进行用户登录验证。
//		server.Use(login.CheckLogin())*/
//	//两个中间件：第一个是用来提取session的，第二个是用来登陆校验的

//	//  单机单实例部署 考虑基于内存的memstore实现，多实例部署，redis
//	//memstore 是基于内存实现的
//	//参数一是authentication key 身份验证
//	//参数二encryption key 数据加密
//	// + 授权(权限控制)就是信息安全的三个核心概念
//	//最好64或者32位
//	//百度 -》 生成密码 -》 复制粘贴
//	//store := memstore.NewStore([]byte("05kcS4LEzQcepWhpjjas07YNzJgxL93a"),
//	//	[]byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgX"))
//	//store, err := redis.NewStore(16, "tcp", "localhost:6379",
//	//	"", []byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgX"),
//	//	[]byte("05kcS4LEzQcepWhpjjas07YNzJgxL93a"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	//server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}

//func initDB() *gorm.DB {
//	dsn := "root:root@tcp(localhost:13316)/webook"
//	db, err := gorm.Open(mysql.Open(dsn))
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTables(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
