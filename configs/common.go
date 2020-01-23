package configs

import (
	"errors"
	"fmt"
	"strings"
)

// EnvPrefix .
const EnvPrefix = "AUTH_"

// EnvError .
type EnvError struct {
	Func string
	Err  error
}

func (e *EnvError) Error() string {
	return "configs." + e.Func + ": " + e.Err.Error()
}

func missingRequirementError(fn string, missed []string) *EnvError {
	const errMessage = "must set '%s' environment variable"
	err := fmt.Sprintf(errMessage, strings.Join(missed, ", "))
	return &EnvError{fn, errors.New(err)}
}
