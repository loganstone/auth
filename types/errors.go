package types

// ErroCodes .
const (
	UnknownError = 1000 + iota
	ValidateError
	DBTransactionError
	MarshalJSONError
)

// ErroCodes .
const (
	NotFoundUser = 2000 + iota
	IncorrectPassword
	UserAlreadyExists
)

// Error .
type Error struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
