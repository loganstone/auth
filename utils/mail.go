package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
)

const (
	localHost = "127.0.0.1"
	port      = "25"
)

// SendMail 는 local postfix 로 email 을 보낸다.
func SendMail(fromName, fromEmail, toEmail, subject, body string) {
	c, err := smtp.Dial(localHost + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	toHeader := toEmail
	from := mail.Address{
		Name:    fromName,
		Address: fromEmail,
	}
	fromHeader := from.String()
	subjectHeader := subject
	header := make(map[string]string)
	header["To"] = toHeader
	header["From"] = fromHeader
	header["Subject"] = subjectHeader
	header["Content-Type"] = `text/html; charset="UTF-8"`
	msg := ""

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	msg += "\r\n" + body

	c.Mail(fromEmail)
	c.Rcpt(toEmail)

	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(msg)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}
