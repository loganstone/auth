package payload

// Internal server error codes.
const (
	ErrorCodeUnknown = 1000 + iota
	ErrorCodeDBTransaction
	ErrorCodeMarshalJSON
	ErrorCodeSendEmail
)

// Bind erro codes.
const (
	ErrorCodeBindURI = 2000 + iota
	ErrorCodeBindJSON
)

// User error codes.
const (
	ErrorCodeNotFoundUser = 3000 + iota
	ErrorCodeUserAlreadyExists
	ErrorCodeSetPassword
	ErrorCodeIncorrectPassword
)
