package web

import (
	"errors"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	EmailReGexPattern    = "/^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$/"
	PasswordReGexPattern = "^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{1,9}$"
)

/*
UserHandler 所有与用户有关的路由定义在这个Handler上
RegisterRouter方法 用来注册路由
*/
type UserHandler struct {
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	svc              *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(EmailReGexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(PasswordReGexPattern, regexp.None),
		svc:              svc,
	}
}

func (h *UserHandler) RegisterRouter(server *gin.Engine) {
	/*
		RESTfull 风格
		service.POST("/user.go",h.SignUp)
		service.PUT("/user.go",h.SignUp)
		service.GET("/user.go/:id",h.Profile)
						 /:username

	*/

	/*	//POST方法 把前端数据推给后端
		service.POST("/users/signup", h.SignUp)
		service.POST("/users/login", h.Login)
		service.POST("/users/edit", h.Edit)
		service.GET("/users/profile", h.Profile)
	*/

	//分组路由
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJST)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	//• 接收请求并校验 • 调用业务逻辑处理请求 • 根据业务逻辑处理结果返回响应
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	println("------------")
	println(req.Email, req.Password)
	/*
		// regex.Match()
		isEmail, err := regexp.Match(EmailReGexPattern, []byte(req.Email))
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		if isEmail {
			ctx.String(http.StatusOK, "邮箱格式错误")
			return
		}*/

	isEmail, err := h.emailRegexExp.MatchString(EmailReGexPattern)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if isEmail {
		println(isEmail)
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	isPass, err := h.emailRegexExp.MatchString(PasswordReGexPattern)
	if err != nil {
		println(err)
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if isPass {
		ctx.String(http.StatusOK, "密码格式错误，至少包含字母、数字、特殊字符，1-9位")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不同")
		return
	}

	err = h.svc.SingUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "hello,正在登录...")
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.String(http.StatusOK, "邮箱冲突,请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}

}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			//十五分钟
			MaxAge: 900,
		})
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

// Edit 修改
func (h *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "edit")
}

// Profile 拿到用户基本信息
func (h *UserHandler) Profile(ctx *gin.Context) {

}

func (h *UserHandler) LoginJST(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch {
	case err == nil:
		uc := UserClaims{
			Uid: u.Id,
			RegisteredClaims: jwt.RegisteredClaims{
				// 30分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		// token 是结构体，改成jwt字节切片传给前端
		tokenString, err := token.SignedString(JWTKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
		}
		ctx.Header("x-jwt-token", tokenString)
		ctx.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

var JWTKey = []byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgX")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
