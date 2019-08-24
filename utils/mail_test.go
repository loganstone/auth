package utils

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

const (
	name    = "Johndoe"
	from    = "johndoe@mail.com"
	to      = "janedoe@mail.com"
	subject = "For sale"
	body    = "Baby shoes. Never worn. "
)

func TestMakeTextHTMLHeader(t *testing.T) {
	expected := map[string]string{
		"Content-Type": `text/html; charset="UTF-8"`,
		"From":         fmt.Sprintf(`"%s" <%s>`, name, from),
		"To":           to,
		"Subject":      subject,
	}
	email := NewEmail(
		name, from, to, subject, body)
	email.makeTextHTMLHeader()

	assert.Equal(t, email.header, expected)
}
