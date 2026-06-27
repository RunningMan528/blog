package middleware

import (
	"blog/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(ctx, "Authorization header is required")
			ctx.Abort()
			return
		}
		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.Unauthorized(ctx, "Bearer token is required")
			ctx.Abort()
			return
		}
		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Unauthorized(ctx, "Invalid token")
			ctx.Abort()
			return
		}
		// 将用户信息存储到上下文中
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
