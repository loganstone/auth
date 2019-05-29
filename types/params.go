package types

// AddUserParams ..
type AddUserParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// AuthenticateParams ..
type AuthenticateParams struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
