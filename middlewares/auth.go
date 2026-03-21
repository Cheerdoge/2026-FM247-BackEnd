package middleware

import (
	handler "2026-FM247-BackEnd/handlers"
	"2026-FM247-BackEnd/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenblacklistservice handler.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.FailWithMessage(c, "请先登录")
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		claims, err := utils.ValidateToken(token)
		if err != nil {
			handler.FailWithMessage(c, "token无效或已过期")
			c.Abort()
			return
		}

		ok, err := tokenblacklistservice.IsBlacklisted(claims.Jti)
		if err != nil {
			handler.FailWithMessage(c, "服务器内部错误")
			c.Abort()
			return
		}
		if ok {
			handler.FailWithMessage(c, "token已被注销,请重新登录")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("jti", claims.Jti)
		c.Next()
	}
}
