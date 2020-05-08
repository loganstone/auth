package configs

import (
	"fmt"
	"net"
	"net/smtp"
	"strconv"
)

const (
	defaultSMTPHost = "127.0.0.1"
	defaultSMTPPort = 25
)

var smtpPort = defaultSMTPPort

// SMTPConfig .
type SMTPConfig struct {
	host string
	port int
}

// Addr .
func (c *SMTPConfig) Addr() string {
	return net.JoinHostPort(c.host, strconv.Itoa(c.port))
}

// DialAndQuit .
func (c *SMTPConfig) DialAndQuit() error {
	con, err := smtp.Dial(c.Addr())
	if err != nil {
		return fmt.Errorf("smtp server dial: %w", err)
	}
	defer con.Quit()
	return nil
}

// SetSMTPPort .
func SetSMTPPort(port int) {
	smtpPort = port
}

// SMTP .
func SMTP() *SMTPConfig {
	return &SMTPConfig{
		defaultSMTPHost,
		smtpPort,
	}
}
