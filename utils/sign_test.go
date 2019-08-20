package utils

import (
	"encoding/json"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

type Data struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func TestSignAndLoad(t *testing.T) {
	data := Data{
		Email: "test@email.com",
		Name:  "Logan",
	}

	val, err := json.Marshal(data)
	assert.Equal(t, err, nil)

	signed, err := Sign(val)
	assert.Equal(t, err, nil)

	loaded := Load(signed)

	assert.Equal(t, val, loaded)
}
