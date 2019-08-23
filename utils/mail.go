package utils

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"
)

const (
	localHost = "127.0.0.1"
	port      = "25"
)

// Email .
type Email struct {
	name    string
	from    string
	to      string
	subject string
	body    string

	header  map[string]string
	message string
}

// NewEmail .
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

// SendToLocalPostfix 는 local postfix 로 email 을 보낸다.
func (m *Email) SendToLocalPostfix() error {
	c, err := smtp.Dial(localHost + ":" + port)
	if err != nil {
		return err
	}
	defer c.Close()

	c.Mail(m.from)
	c.Rcpt(m.to)

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	m.makeTextHTMLHeader()
	m.makeMessage()

	buf := bytes.NewBufferString(m.message)
	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}

	return nil
}
