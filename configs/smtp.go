package configs

import (
	"fmt"
	"net"
	"net/smtp"
)

const (
	defaultSMTPHost = "127.0.0.1"
	defaultSMTPPort = "25"
)

// SMTPConfigs .
type SMTPConfigs struct {
	host string
	port string
}

// Addr .
func (c *SMTPConfigs) Addr() string {
	return net.JoinHostPort(c.host, c.port)
}

// DialAndQuit .
func (c *SMTPConfigs) DialAndQuit() error {
	con, err := smtp.Dial(c.Addr())
	if err != nil {
		return fmt.Errorf("smtp server dial: %w", err)
	}
	defer con.Quit()
	return nil
}

// SMTP .
func SMTP() *SMTPConfigs {
	return &SMTPConfigs{
		defaultSMTPHost,
		defaultSMTPPort,
	}
}
