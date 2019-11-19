package utils

import (
	"bufio"
	"fmt"
	"io"
	"net"
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
func TestSendToLocalPostfix(t *testing.T) {
	ln := newLocalListener(t)
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
		if err := serverHandle(c, t); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()

	go func() {
		defer close(clientDone)
		email := NewEmail(name, from, to, subject, body)
		err := email.sendToLocalPostfix(ln.Addr().String())
		assert.Nil(t, err)
	}()

	<-clientDone
	<-serverDone
}

func TestNameFromEmail(t *testing.T) {
	expected := "johndoe"
	name := NameFromEmail(from)
	assert.Equal(t, expected, name)
}

// reference - https://golang.org/src/net/smtp/smtp_test.go
func newLocalListener(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", net.JoinHostPort(localHost, testSMTPPort))
	if err != nil {
		ln, err = net.Listen("tcp6", net.JoinHostPort("::1", testSMTPPort))
	}
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

// reference - https://golang.org/src/net/smtp/smtp_test.go
type smtpSender struct {
	w io.Writer
}

// reference - https://golang.org/src/net/smtp/smtp_test.go
func (s smtpSender) send(f string) {
	s.w.Write([]byte(f + "\r\n"))
}

// reference - https://golang.org/src/net/smtp/smtp_test.go
func serverHandle(c net.Conn, t *testing.T) error {
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
			send("354 send the mail data, end with .")
			send("250 Ok")
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
