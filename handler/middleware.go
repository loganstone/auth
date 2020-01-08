package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// Authorize .
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := configs.App()
		dbConf := configs.DB()
		con := db.Connection(dbConf.ConnectionString(), dbConf.Echo)
		defer con.Close()

		reqToken := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(reqToken, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sessionToken := bearerToken[1]
		sessionClaims, err := utils.ParseSessionJWT(
			sessionToken, conf.JWTSigninKey)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := db.User{}

		if con.First(&user, sessionClaims.UserID).RecordNotFound() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Email != sessionClaims.UserEmail {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("AuthorizedUser", user)
		c.Next()
	}
}

// AuthorizedUserIsAdmin .
func AuthorizedUserIsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizedUser, err := AuthorizedUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				payload.ErrorAuthorizedUser(err))
			return
		}

		if !authorizedUser.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("AuthorizedUserIsAdmin", true)
		c.Next()
	}
}

// RequesterIsAuthorizedUser .
func RequesterIsAuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		if email == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		authorizedUser, err := AuthorizedUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				payload.ErrorAuthorizedUser(err))
			return
		}

		if authorizedUser.Email != email {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("RequesterIsAuthorizedUser", true)
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

// Ref - https://sourcegraph.com/github.com/gin-gonic/gin/-/blob/logger.go#L131
var logFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	requestID := param.Request.Header.Get("Request-ID")
	return fmt.Sprintf("[REQUEST ID - %s] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
		requestID,
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// LogFormat .
func LogFormat() gin.HandlerFunc {
	return gin.LoggerWithFormatter(logFormatter)
}
