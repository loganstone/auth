package types

// ErroCodes .
const (
	UnknownError = 1000 + iota
	NotFoundUser
	ValidateError
)

// Error .
type Error struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
