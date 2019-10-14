package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/utils"
)

// Authorize .
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(reqToken, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionToken := bearerToken[1]
		sessionClaims, err := utils.ParseJWTSessionToken(sessionToken)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		con := db.Connection()
		defer con.Close()

		user := models.User{}

		if con.First(&user, sessionClaims.UserID).RecordNotFound() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Email != sessionClaims.UserEmail {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("SessionUser", user)
		c.Next()
	}
}

// RequestID .
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Set("Request-ID", uuid.New().String())
		c.Next()
	}
}
