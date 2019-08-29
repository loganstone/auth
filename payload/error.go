package payload

import (
	"github.com/gin-gonic/gin"
)

// ErrorWithCode .
func ErrorWithCode(code int, message string) gin.H {
	return gin.H{
		"error_code":    code,
		"error_message": message,
	}
}

// ErrorDBTransaction .
func ErrorDBTransaction(message string) gin.H {
	return ErrorWithCode(ErrorCodeDBTransaction, message)
}

// ErrorBindURI .
func ErrorBindURI(message string) gin.H {
	return ErrorWithCode(ErrorCodeBindURI, message)
}

// ErrorBindJSON .
func ErrorBindJSON(message string) gin.H {
	return ErrorWithCode(ErrorCodeBindJSON, message)
}

// NotFoundUser .
func NotFoundUser() gin.H {
	return ErrorWithCode(ErrorCodeNotFoundUser, "not such user")
}

// UserAlreadyExists .
func UserAlreadyExists() gin.H {
	return ErrorWithCode(
		ErrorCodeUserAlreadyExists, "user already exists")
}

// ErrorSetPassword .
func ErrorSetPassword(message string) gin.H {
	return ErrorWithCode(ErrorCodeSetPassword, message)
}
