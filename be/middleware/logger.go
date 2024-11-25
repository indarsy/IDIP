package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 日志中间件
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装ResponseWriter以记录响应体
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 记录日志
		logger.Info("HTTP Request",
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("request", string(requestBody)),
			zap.String("response", blw.body.String()),
		)
	}
}
