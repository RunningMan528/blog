package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ===== 响应结构体 =====
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(ctx *gin.Context, message string) {
	Error(ctx, http.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(ctx *gin.Context, message string) {
	Error(ctx, http.StatusUnauthorized, message)
}

// Forbidden 403错误
func Forbidden(ctx *gin.Context, message string) {
	Error(ctx, http.StatusForbidden, message)
}

// NotFound 404错误
func NotFound(ctx *gin.Context, message string) {
	Error(ctx, http.StatusNotFound, message)
}

// InternalServerError 500错误
func InternalServerError(ctx *gin.Context, message string) {
	Error(ctx, http.StatusInternalServerError, message)
}
