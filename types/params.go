package types

// AddUserParams ..
type AddUserParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// SigninParams .
type SigninParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
