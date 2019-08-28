package response

import (
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/types"
)

// ErrorCode .
func ErrorCode(code int, message string) gin.H {
	return gin.H{
		"error_code":    code,
		"error_message": message,
	}
}

// DBTransactionError .
func DBTransactionError(message string) gin.H {
	return ErrorCode(types.DBTransactionError, message)
}

// BindURIError .
func BindURIError(message string) gin.H {
	return ErrorCode(types.BindURIError, message)
}

// BindJSONError .
func BindJSONError(message string) gin.H {
	return ErrorCode(types.BindJSONError, message)
}

// NotFoundUser .
func NotFoundUser() gin.H {
	return ErrorCode(types.NotFoundUser, "not such user")
}

// UserAlreadyExists .
func UserAlreadyExists() gin.H {
	return ErrorCode(types.UserAlreadyExists, "user already exists")
}

// SetPasswordError .
func SetPasswordError(message string) gin.H {
	return ErrorCode(types.SetPasswordError, message)
}
