package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionWithBadDsn(t *testing.T) {
	baddsn := "baddsn"
	_, err := Connection(baddsn, true)
	expectedError := "invalid DSN: missing the slash separating the database name"
	assert.EqualError(t, err, expectedError)
}
