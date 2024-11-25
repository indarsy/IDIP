package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery中间件
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := debug.Stack()

				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(stack)),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": fmt.Sprintf("Internal Server Error: %v", err),
				})
			}
		}()

		c.Next()
	}
}
