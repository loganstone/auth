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
	"github.com/loganstone/auth/utils"
)

// Authorize .
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := configs.App()
		con := DBConnOrAbort(c)
		if con == nil {
			return
		}

		reqToken := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(reqToken, " ")
		if len(bearerToken) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseSessionJWT(bearerToken[1], conf.JWTSigninKey)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := db.User{}
		if con.First(&user, claims.UserID).RecordNotFound() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Email != claims.UserEmail {
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
		user, err := AuthorizedUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				NewErrResWithErr(ErrorCodeAuthorizedUser, err))
			return
		}

		if !user.IsAdmin {
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

		user, err := AuthorizedUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				NewErrResWithErr(ErrorCodeAuthorizedUser, err))
			return
		}

		if user.Email != email {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("RequesterIsAuthorizedUser", true)
		c.Next()
	}
}

// DBConnection .
func DBConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConf, err := configs.DB()
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				NewErrResWithErr(ErrorCodeDBEnv, err))
		}
		dbConf.SetMode(gin.Mode())
		con, err := db.Connection(dbConf.ConnectionString(), dbConf.Echo)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				NewErrResWithErr(ErrorCodeDBConn, err))
		}
		defer con.Close()
		c.Set("DBConnection", con)
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
