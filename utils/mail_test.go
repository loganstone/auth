package utils

import (
	"fmt"
	"strings"
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

func TestMakeMessage(t *testing.T) {
	expected := true
	email := NewEmail(
		name, from, to, subject, body)
	email.makeTextHTMLHeader()
	email.makeMessage()

	contained := strings.Contains(
		email.message, fmt.Sprintf(`From: "%s" <%s>`, name, from))
	assert.Equal(t, contained, expected)

	contained = strings.Contains(
		email.message, fmt.Sprintf(`To: %s`, to))
	assert.Equal(t, contained, expected)

	contained = strings.Contains(
		email.message, fmt.Sprintf(`Subject: %s`, subject))
	assert.Equal(t, contained, expected)

	contained = strings.Contains(
		email.message, `Content-Type: text/html; charset="UTF-8"`)
	assert.Equal(t, contained, expected)

	contained = strings.Contains(email.message, "\n")
	assert.Equal(t, contained, expected)

	contained = strings.Contains(email.message, body)
	assert.Equal(t, contained, expected)
}
