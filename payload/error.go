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

// ErrorMarshalJSON .
func ErrorMarshalJSON(message string) gin.H {
	return ErrorWithCode(ErrorCodeMarshalJSON, message)
}

// ErrorUnMarshalJSON .
func ErrorUnMarshalJSON(message string) gin.H {
	return ErrorWithCode(ErrorCodeUnMarshalJSON, message)
}

// ErrorSignJWTToken .
func ErrorSignJWTToken(message string) gin.H {
	return ErrorWithCode(ErrorCodeSignJWTToken, message)
}

// ErrorParseJWTToken .
func ErrorParseJWTToken(message string) gin.H {
	return ErrorWithCode(ErrorCodeParseJWTToken, message)
}

// ErrorSendEmail .
func ErrorSendEmail(message string) gin.H {
	return ErrorWithCode(ErrorCodeSendEmail, message)
}

// ErrorTmplExecute .
func ErrorTmplExecute(message string) gin.H {
	return ErrorWithCode(ErrorCodeTmplExecute, message)
}

// ErrorBindURI .
func ErrorBindURI(message string) gin.H {
	return ErrorWithCode(ErrorCodeBindURI, message)
}

// ErrorBindJSON .
func ErrorBindJSON(message string) gin.H {
	return ErrorWithCode(ErrorCodeBindJSON, message)
}

// ErrorBadPage .
func ErrorBadPage(message string) gin.H {
	return ErrorWithCode(ErrorCodeBadPage, message)
}

// ErrorBadPageSize .
func ErrorBadPageSize(message string) gin.H {
	return ErrorWithCode(ErrorCodeBadPageSize, message)
}

// ErrorExpiredToken .
func ErrorExpiredToken() gin.H {
	return ErrorWithCode(ErrorCodeExpiredToken, "expired token")
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

// ErrorSession .
func ErrorSession(err error) gin.H {
	return ErrorWithCode(ErrorCodeWrongSession, err.Error())
}
