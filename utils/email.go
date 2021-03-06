package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"strings"
)

// Email sends an email.
type Email struct {
	name    string
	from    string
	to      string
	subject string
	body    string

	header  map[string]string
	message string
	wc      io.WriteCloser
}

// NewEmail creates and returns a new Email object.
func NewEmail(name, from, to, subject, body string) *Email {
	return &Email{
		name:    name,
		from:    from,
		to:      to,
		subject: subject,
		body:    body,

		header: map[string]string{},
	}
}

func (m *Email) makeHeader(contentType string) {
	from := mail.Address{
		Name:    m.name,
		Address: m.from,
	}
	m.header["To"] = m.to
	m.header["From"] = from.String()
	m.header["Subject"] = m.subject
	m.header["Content-Type"] = contentType
}

func (m *Email) makeTextHTMLHeader() {
	m.makeHeader(`text/html; charset="UTF-8"`)
}

func (m *Email) makeMessage() {
	for k, v := range m.header {
		m.message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	m.message += "\r\n" + m.body
}

// Send 는 local postfix 로 email 을 보낸다.
func (m *Email) Send(addr string) error {
	return m.sendToSMTPServer(addr)
}

func (m *Email) sendToSMTPServer(addr string) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Quit()

	c.Mail(m.from)
	c.Rcpt(m.to)

	if m.wc == nil {
		if m.wc, err = c.Data(); err != nil {
			return err
		}
		defer m.wc.Close()
	}

	m.makeTextHTMLHeader()
	m.makeMessage()

	buf := bytes.NewBufferString(m.message)
	if _, err = buf.WriteTo(m.wc); err != nil {
		return err
	}

	return nil
}

// NameFromEmail .
func NameFromEmail(email string) string {
	return strings.Split(email, "@")[0]
}
