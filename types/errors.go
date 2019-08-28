package types

// ErroCodes .
const (
	UnknownError = 1000 + iota
	DBTransactionError
	MarshalJSONError
	SendEmailError
)

// ErroCodes .
const (
	BindURIError = 2000 + iota
	BindJSONError
)

// ErroCodes .
const (
	NotFoundUser = 3000 + iota
	UserAlreadyExists
	SetPasswordError
	IncorrectPassword
)
