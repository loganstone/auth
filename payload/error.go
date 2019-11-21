package payload

// ErrorCodeResponse .
type ErrorCodeResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// ErrorResponse .
func ErrorResponse(code int, message string) ErrorCodeResponse {
	return ErrorCodeResponse{code, message}
}

// ErrorDBTransaction .
func ErrorDBTransaction(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeDBTransaction, message)
}

// ErrorMarshalJSON .
func ErrorMarshalJSON(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeMarshalJSON, message)
}

// ErrorUnMarshalJSON .
func ErrorUnMarshalJSON(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeUnMarshalJSON, message)
}

// ErrorSignJWTToken .
func ErrorSignJWTToken(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeSignJWTToken, message)
}

// ErrorParseJWTToken .
func ErrorParseJWTToken(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeParseJWTToken, message)
}

// ErrorSendEmail .
func ErrorSendEmail(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeSendEmail, message)
}

// ErrorTmplExecute .
func ErrorTmplExecute(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeTmplExecute, message)
}

// ErrorBindURI .
func ErrorBindURI(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeBindURI, message)
}

// ErrorBindJSON .
func ErrorBindJSON(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeBindJSON, message)
}

// ErrorBadPage .
func ErrorBadPage(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeBadPage, message)
}

// ErrorBadPageSize .
func ErrorBadPageSize(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeBadPageSize, message)
}

// ErrorExpiredToken .
func ErrorExpiredToken() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeExpiredToken, "expired token")
}

// NotFoundUser .
func NotFoundUser() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeNotFoundUser, "not such user")
}

// UserAlreadyExists .
func UserAlreadyExists() ErrorCodeResponse {
	return ErrorResponse(
		ErrorCodeUserAlreadyExists, "user already exists")
}

// ErrorSetPassword .
func ErrorSetPassword(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeSetPassword, message)
}

// ErrorSession .
func ErrorSession(err error) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeWrongSession, err.Error())
}

// ErrorIncorrectOTP .
func ErrorIncorrectOTP() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeIncorrectOTP, "OTP is Incorrect")
}

// ErrorOTPAlreadyRegistered .
func ErrorOTPAlreadyRegistered() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeOTPAlreadyRegistered, "OTP has already been registered")
}

// ErrorOTPNotRegistered .
func ErrorOTPNotRegistered() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeOTPNotRegistered, "OTP not registered")
}

// ErrorEmptyOTPSecretKey .
func ErrorEmptyOTPSecretKey() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeEmptyOTPSecretKey, "empty OTP secert key")
}

// ErrorEmptyOTPBackupCodes .
func ErrorEmptyOTPBackupCodes(message string) ErrorCodeResponse {
	return ErrorResponse(ErrorCodeEmptyOTPBackupCodes, message)
}

// ErrorRequireVerifyOTP .
func ErrorRequireVerifyOTP() ErrorCodeResponse {
	return ErrorResponse(ErrorCodeRequireVerifyOTP, "required verify OTP")
}
