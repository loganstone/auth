package configs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvError(t *testing.T) {
	const fnTest = "Test"
	const errMessage = "test error"
	expected := "configs." + fnTest + ": " + errMessage
	err := &EnvError{fnTest, errors.New(errMessage)}
	assert.Equal(t, expected, err.Error())
}
