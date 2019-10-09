package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/loganstone/auth/utils"
)

// Authorize .
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(reqToken, " ")
		if len(bearerToken) != 2 {
			// TODO(hs.lee): error_code 처리
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionToken := bearerToken[1]
		sessionClaims, err := utils.ParseJWTSessionToken(sessionToken)
		if err != nil {
			// TODO(hs.lee): error_code 처리
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("SessionUser", sessionClaims.SessionUser)
		c.Next()
	}
}

// RequestID .
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Request-Id", uuid.New().String())
		c.Next()
	}
}
