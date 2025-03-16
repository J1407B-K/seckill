package utils

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

// JWTAuthMiddleware 是一个用于验证 JWT 的中间件
func JWTAuthMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从 Header 中获取 Authorization 字段（token）
		authHeader := ctx.GetHeader("Authorization")
		if string(authHeader) == "" {
			ctx.JSON(http.StatusUnauthorized, utils.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		// 检查 Bearer 格式（我也不知道是啥，似乎就要这么写，在Authorization的[0]是Bearer，[1]是token）
		tokenString := string(authHeader)
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析
		claims, err := ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// 将解析后的用户信息存入上下文
		ctx.Set("username", claims.Username)

		// 继续处理请求
		ctx.Next(c)
	}
}
