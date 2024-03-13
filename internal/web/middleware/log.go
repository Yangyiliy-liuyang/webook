package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type LogMiddlewareBuilder struct {
	LogFun        func(ctx context.Context, l AcessLog)
	allowReqBody  bool
	allowRespBody bool
}

type AcessLog struct {
	Path     string        `json:"path"`
	Method   string        `json:"method"`
	Status   int           `json:"status"`
	ReqBody  string        `json:"req_body"`
	RespBody string        `json:"resp_body"`
	Duration time.Duration `json:"duration"`
}

func NewLogMiddlewareBuilder(logFun func(ctx context.Context, l AcessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		LogFun: logFun,
	}
}

func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	l.allowReqBody = true
	return l
}

func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	l.allowRespBody = true
	return l
}

func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if len(path) > 1024 {
			path = path[:1024]
		}
		method := ctx.Request.Method
		al := AcessLog{
			Path:   path,
			Method: method,
		}
		if l.allowReqBody {
			// Request.Body 是一个Stream 只能读取一次，需要重新赋值
			body, _ := ctx.GetRawData()
			al.ReqBody = string(body)
			//ctx.Request.Body = io.NopCloser(bytes.NewReader(body))// 重置请求体
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置请求体
		}
		start := time.Now()
		if l.allowRespBody {
			w := &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
			ctx.Writer = w
		}
		// 记录开始时间
		defer func() {
			// 记录结束时间
			al.Duration = time.Since(start)
			//al.Duration = time.Now().Sub(start)
			l.LogFun(ctx, al)
		}()
		// 执行下一个中间件
		ctx.Next()

	}
}

type responseWriter struct {
	gin.ResponseWriter
	al *AcessLog
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write([]byte(w.al.RespBody))
}
func (w *responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
