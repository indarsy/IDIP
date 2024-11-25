package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应码定义
const (
	SUCCESS = 0
	ERROR   = 1
)

// 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SUCCESS,
		Message: "success",
		Data:    data,
	})
}

// 错误响应
func Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:    ERROR,
		Message: err.Error(),
		Data:    nil,
	})
}

// 自定义状态码响应
func ResponseWithCode(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
