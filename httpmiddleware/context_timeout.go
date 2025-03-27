package middleware

import (
	"context"
	"github.com/niclausse/webkit/zlog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	contextTimeout = time.Second * 15
)

func ContextTimeout(timeout time.Duration) gin.HandlerFunc {
	if timeout == 0 {
		timeout = contextTimeout
	}
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// replace original context with timeout context
		c.Request = c.Request.WithContext(ctx)

		// job done flag
		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()           // execute next handler chain(s)
			done <- struct{}{} // set flag when all handlers were finished
		}()

		select {
		case r := <-panicChan:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errNo":  http.StatusInternalServerError,
				"errMsg": "服务器内部错误",
			})
			zlog.Errorf("[http] panic: %v", r)
			return
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"errNo":  http.StatusRequestTimeout,
				"errMsg": "请求超时",
			})
			zlog.Errorf("[http] context deadline, %s - %s - %s", c.Request.URL, c.Request.RemoteAddr, c.Request.UserAgent())
			return
		case <-done:
			return
		}
	}
}
