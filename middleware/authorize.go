package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/utils"
)

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
