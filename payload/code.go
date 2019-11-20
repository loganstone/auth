package payload

// Internal server error codes.
const (
	ErrorCodeUnknown = iota + 1000
	ErrorCodeDBTransaction
	ErrorCodeMarshalJSON
	ErrorCodeUnMarshalJSON
	ErrorCodeSignJWTToken
	ErrorCodeParseJWTToken
	ErrorCodeSendEmail
	ErrorCodeTmplExecute
)

// Parameter error codes.
const (
	ErrorCodeBindURI = iota + 2000
	ErrorCodeBindJSON
	ErrorCodeBadPage
	ErrorCodeBadPageSize
	ErrorCodeExpiredToken
)

// User error codes.
const (
	ErrorCodeNotFoundUser = iota + 3000
	ErrorCodeUserAlreadyExists
	ErrorCodeSetPassword
	ErrorCodeIncorrectPassword
	ErrorCodeOTPProvisioningURI
)

// Session error codes.
const (
	ErrorCodeWrongSession = iota + 4000
)
