package MiddleWare

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type LogMiddlewareBuilder struct {
	logFn         func(ctx context.Context, l AccessLog)
	allowReqBody  bool
	allowRespBody bool
	PathThreshold int
	BodyThreshold int
}

func NewLogMiddlewareBuilder(logFn func(ctx context.Context, l AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFn:         logFn,
		PathThreshold: 1024,
		BodyThreshold: 2048,
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
		if len(path) > l.PathThreshold {
			path = path[:l.PathThreshold]
		}
		al := AccessLog{
			Path:   path,
			Method: ctx.Request.Method,
		}
		if l.allowReqBody {
			body, _ := ctx.GetRawData()
			if len(body) > l.BodyThreshold {
				al.ReqBody = string(body[:l.BodyThreshold])
			} else {
				al.ReqBody = string(body)
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		start := time.Now()
		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}
		defer func() {
			al.Duration = time.Since(start)
			l.logFn(ctx, al)
		}()
		// 拿到响应
		// 执行下一个中间件及业务逻辑，直到响应
		ctx.Next()
	}
}

type AccessLog struct {
	Path     string        `json:"path"`
	Method   string        `json:"method"`
	ReqBody  string        `json:"req_body"`
	Status   int           `json:"status"`
	RespBody string        `json:"resp_body"`
	Duration time.Duration `json:"duration"`
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func (r *responseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.al.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
