package payload

// Internal server error codes.
const (
	ErrorCodeUnknown = 1000 + iota
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
	ErrorCodeBindURI = 2000 + iota
	ErrorCodeBindJSON
	ErrorCodeBadPage
	ErrorCodeBadPageSize
	ErrorCodeExpiredToken
)

// User error codes.
const (
	ErrorCodeNotFoundUser = 3000 + iota
	ErrorCodeUserAlreadyExists
	ErrorCodeSetPassword
	ErrorCodeIncorrectPassword
)
