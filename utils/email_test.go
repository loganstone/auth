package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	name    = "Johndoe"
	from    = "johndoe@mail.com"
	to      = "janedoe@mail.com"
	subject = "For sale"
	body    = "Baby shoes. Never worn. "
)

// reference - https://golang.org/src/net/smtp/smtp.go
// type dataCloser struct
type badCloser struct {
	c *smtp.Client
	io.WriteCloser
}

func (d *badCloser) Close() error {
	d.WriteCloser.Close()
	_, _, err := d.c.Text.ReadResponse(250)
	return err
}

func (d *badCloser) Write(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

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

// reference - https://golang.org/src/net/smtp/smtp_test.go
func TestSendToSMTP(t *testing.T) {
	ln, err := NewLocalListener(MockSMTPPort)
	assert.NoError(t, err)
	defer ln.Close()

	clientDone := make(chan bool)
	serverDone := make(chan bool)

	go func() {
		defer close(serverDone)
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("server accept: %v", err)
			return
		}
		defer c.Close()
		handler := MockSMTPHandler{
			Con: c, Name: name, From: from, To: to,
			Subject: subject, Body: body,
		}
		if err := handler.Handle(); err != nil {
			t.Errorf("mock smtp handle error: %v", err)
		}
	}()

	go func() {
		defer close(clientDone)
		email := NewEmail(name, from, to, subject, body)
		err := email.Send(ln.Addr().String())
		assert.NoError(t, err)
	}()

	<-clientDone
	<-serverDone
}

func TestSendWithBadSMTPServer(t *testing.T) {
	expectedError := errors.New("dial tcp: address bad smtp address: missing port in address")
	email := NewEmail(name, from, to, subject, body)
	err := email.Send("bad smtp address")
	assert.EqualError(t, expectedError, fmt.Sprint(err))
}

func TestSendWithBadServerHandleForData(t *testing.T) {
	ln, err := NewLocalListener(MockSMTPPort)
	assert.NoError(t, err)
	defer ln.Close()

	expectedError := errors.New("EOF")
	clientDone := make(chan bool)
	serverDone := make(chan bool)

	go func() {
		defer close(serverDone)
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("server accept: %v", err)
			return
		}
		defer c.Close()
		if err := badServerHandleForData(c, t); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()

	go func() {
		defer close(clientDone)
		email := NewEmail(name, from, to, subject, body)
		err := email.Send(ln.Addr().String())
		assert.EqualError(t, expectedError, fmt.Sprint(err))
	}()

	<-clientDone
	<-serverDone
}

func TestSendWithBadCloser(t *testing.T) {
	ln, err := NewLocalListener(MockSMTPPort)
	assert.NoError(t, err)
	defer ln.Close()

	expectedError := "error"
	clientDone := make(chan bool)
	serverDone := make(chan bool)

	go func() {
		defer close(serverDone)
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("server accept: %v", err)
			return
		}
		defer c.Close()
		handler := MockSMTPHandler{
			Con: c, Name: name, From: from, To: to,
			Subject: subject, Body: body,
		}
		if err := handler.Handle(); err != nil {
			t.Errorf("mock smtp handle error: %v", err)
		}
	}()

	go func() {
		defer close(clientDone)
		email := NewEmail(name, from, to, subject, body)
		email.wc = &badCloser{}
		err := email.Send(ln.Addr().String())
		assert.EqualError(t, err, expectedError)
	}()

	<-clientDone
	<-serverDone
}

func TestNameFromEmail(t *testing.T) {
	expected := "johndoe"
	name := NameFromEmail(from)
	assert.Equal(t, expected, name)
}

func badServerHandleForData(c net.Conn, t *testing.T) error {
	send := smtpSender{c}.send
	// Important.
	send("220 127.0.0.1 ESMTP service ready")
	s := bufio.NewScanner(c)
	for s.Scan() {
		switch s.Text() {
		case "EHLO localhost":
			send("250 Ok")
		case fmt.Sprintf("MAIL FROM:<%s>", from):
			send("250 Ok")
		case fmt.Sprintf("RCPT TO:<%s>", to):
			send("250 Ok")
		case "DATA":
			return nil
		case fmt.Sprintf("Subject: %s", subject):
		case `Content-Type: text/html; charset="UTF-8"`:
		case fmt.Sprintf(`From: "%s" <%s>`, name, from):
		case fmt.Sprintf("To: %s", to):
		case "":
		case body:
		case ".":
		case "QUIT":
			send("221 127.0.0.1 Service closing transmission channel")
			return nil
		default:
			t.Fatalf("unrecognized command: %q", s.Text())
		}
	}
	return s.Err()
}
