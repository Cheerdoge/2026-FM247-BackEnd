package middleware

import (
	handler "2026-FM247-BackEnd/handlers"
	"2026-FM247-BackEnd/utils"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.GetClaimsFromContext(c)
		if err != nil {
			handler.FailWithMessage(c, "用户信息不存在")
			c.Abort()
			return
		}
		if !claims.IsAdmin {
			handler.FailWithMessage(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
