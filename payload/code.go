package payload

// ErroCodes .
const (
	ErrorCodeUnknown = 1000 + iota
	ErrorCodeDBTransaction
	ErrorCodeMarshalJSON
	ErrorCodeSendEmail
)

// ErroCodes .
const (
	ErrorCodeBindURI = 2000 + iota
	ErrorCodeBindJSON
)

// ErroCodes .
const (
	ErrorCodeNotFoundUser = 3000 + iota
	ErrorCodeUserAlreadyExists
	ErrorCodeSetPassword
	ErrorCodeIncorrectPassword
)
