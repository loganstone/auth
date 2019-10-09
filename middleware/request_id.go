package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Request-Id", uuid.New().String())
		c.Next()
	}
}
