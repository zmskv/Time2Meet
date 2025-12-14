package middleware

import (
	"time2meet/internal/domain/valueobject"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserIDKey = "user_id"
	CtxIPKey     = "ip"
)

func ContextFromHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v := c.GetHeader("X-User-Id"); v != "" {
			if id, err := valueobject.ParseUUID(v); err == nil && id != valueobject.Nil {
				c.Set(CtxUserIDKey, id)
			}
		}
		c.Set(CtxIPKey, c.ClientIP())
		c.Next()
	}
}
