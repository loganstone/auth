package payload

// Internal server error codes.
const (
	ErrorCodeUnknown = iota + 1000
	ErrorCodeDBEnv
	ErrorCodeDBTransaction
	ErrorCodeMarshalJSON
	ErrorCodeUnMarshalJSON
	ErrorCodeSignJWTToken
	ErrorCodeParseJWTToken
	ErrorCodeSendEmail
	ErrorCodeTmplExecute
	ErrorCodeTmplParse
)

// Parameter error codes.
const (
	ErrorCodeBindURI = iota + 2000
	ErrorCodeBindJSON
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
