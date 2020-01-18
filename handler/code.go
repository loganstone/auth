package handler

// Internal server error codes.
const (
	ErrorCodeUnknown = iota + 1000
	ErrorCodeDBEnv
	ErrorCodeDBConn
	ErrorCodeEmptyDBConn
	ErrorCodeWrongDBConn
	ErrorCodeDBTransaction
	ErrorCodeSignJWT
	ErrorCodeParseJWT
	ErrorCodeSendEmail
	ErrorCodeTmplExecute
	ErrorCodeTmplParse
)

// Parameter error codes.
const (
	ErrorCodeBindJSON = iota + 2000
	ErrorCodeBadPage
	ErrorCodeBadPageSize
	ErrorCodeExpiredToken
	ErrorCodeInvalidPassword
)

// User error codes.
const (
	ErrorCodeNotFoundUser = iota + 3000
	ErrorCodeUserAlreadyExists
	ErrorCodeSetPassword
	ErrorCodeIncorrectPassword
	ErrorCodeOTPProvisioningURI
	ErrorCodeIncorrectOTP
	ErrorCodeSetOTPBackupCodes
	ErrorCodeOTPAlreadyRegistered
	ErrorCodeOTPNotRegistered
	ErrorCodeEmptyOTPSecretKey
	ErrorCodeEmptyOTPBackupCodes
	ErrorCodeRequireVerifyOTP
)

// Authorized User error codes.
const (
	ErrorCodeAuthorizedUser = iota + 4000
)
