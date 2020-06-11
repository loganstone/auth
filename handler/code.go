package handler

// Internal server error codes.
const (
	ErrorCodeUnknown = iota + 1000

	ErrorCodeDBEnv
	ErrorCodeDBConn
	ErrorCodeNoDBConn
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

// User data error codes.
const (
	ErrorCodeNotFoundUser = iota + 3000
	ErrorCodeUserAlreadyExists
	ErrorCodeIncorrectPassword
	ErrorCodeSetPassword

	ErrorCodeOTPAlreadyRegistered
	ErrorCodeEmptyOTPSecretKey
	ErrorCodeIncorrectOTP
	ErrorCodeEmptyOTPBackupCodes
	ErrorCodeRequireVerifyOTP

	ErrorCodeOTPProvisioningURI
	ErrorCodeSetOTPBackupCodes
	ErrorCodeOTPNotRegistered
)

// Authorized User error codes.
const (
	ErrorCodeAuthorizedUser = iota + 4000
)
