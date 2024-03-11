package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type JWTHandler struct {
	signingMethod jwt.SigningMethod
	refreshKey    []byte
	cmd           redis.Cmdable
	rcExpiration  time.Duration
}

func newJWTHandler() JWTHandler {
	return JWTHandler{
		signingMethod: jwt.SigningMethodHS512,
		refreshKey:    []byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgC"),
		rcExpiration:  time.Hour * 24 * 7,
	}
}

func (h *JWTHandler) clearToken(ctx *gin.Context) error {
	ctx.Header("x-refresh-token", "")
	ctx.Header("x-jwt-token", "")
	uc := ctx.MustGet("user").(UserClaims)
	err := h.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", uc.Ssid), "", h.rcExpiration).Err()
	return err
}

func (h *JWTHandler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

func (h *JWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	claims := jwt.NewWithClaims(h.signingMethod, rc)
	signedString, err := claims.SignedString(h.refreshKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-refresh-token", signedString)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	// 根据约定 token在Authorization头部
	// Bearer xxx
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		//没有Token
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return authCode
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		//格式错误
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	tokenStr := segs[1]
	return tokenStr
}

var JWTKey = []byte("Cw7kG6rkQi3WUJ7svOrK4KMStXQ6ykgX")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}

func (h *JWTHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 30分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
			Issuer:    "webook",
		}}
	//使用指定的签名方法创建
	token := jwt.NewWithClaims(h.signingMethod, uc)
	// token 是结构体，改成jwt字节切片传给前端
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenString)
	return nil
}
