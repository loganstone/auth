package mocks

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	// SMTPPort .
	SMTPPort = 7777
)

const (
	localHost = "127.0.0.1"
)

// NewLocalListener .
func NewLocalListener(p int) (net.Listener, error) {
	port := strconv.Itoa(p)
	ln, err := net.Listen("tcp", net.JoinHostPort(localHost, port))
	if err != nil {
		ln, err = net.Listen("tcp6", net.JoinHostPort("::1", port))
	}
	if err != nil {
		return nil, err
	}
	return ln, nil
}

type smtpSender struct {
	w io.Writer
}

func (s smtpSender) send(f string) {
	s.w.Write([]byte(f + "\r\n"))
}

// Handler .
type Handler interface {
	Handle() error
}

// SMTPHandler .
type SMTPHandler struct {
	Con     net.Conn
	Name    string
	From    string
	To      string
	Subject string
	Body    string
}

// Handle .
func (h *SMTPHandler) Handle() error {
	send := smtpSender{h.Con}.send
	// Important.
	send("220 127.0.0.1 ESMTP service ready")
	s := bufio.NewScanner(h.Con)
	for s.Scan() {
		txt := s.Text()
		switch {
		case txt == "EHLO localhost":
			send("250 Ok")
		case txt == fmt.Sprintf("MAIL FROM:<%s>", h.From):
			send("250 Ok")
		case txt == fmt.Sprintf("RCPT TO:<%s>", h.To):
			send("250 Ok")
		case txt == "DATA":
			send("354 send the mail data, end with .")
			send("250 Ok")
		case txt == fmt.Sprintf("Subject: %s", h.Subject):
		case txt == `Content-Type: text/html; charset="UTF-8"`:
		case txt == fmt.Sprintf(`From: "%s" <%s>`, h.Name, h.From):
		case txt == fmt.Sprintf("To: %s", h.To):
		case txt == "":
		case txt == ".":
		case strings.Contains(txt, "signup/email/verification"):
		case strings.Contains(h.Body, txt):
		case txt == "QUIT":
			send("221 127.0.0.1 Service closing transmission channel")
			return nil
		default:
			log.Fatalf("unrecognized command: %q", txt)
		}
	}
	return s.Err()
}
