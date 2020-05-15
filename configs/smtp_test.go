package configs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/utils"
)

func TestSMTP(t *testing.T) {
	ln, err := utils.NewLocalListener(utils.MockSMTPPort)
	assert.NoError(t, err)
	defer ln.Close()

	go func() {
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("local listener accept: %v", err)
			return
		}
		defer c.Close()
		handler := utils.MockSMTPHandler{Con: c}
		if err := handler.Handle(); err != nil {
			t.Errorf("mock smtp handle error: %v", err)
		}
	}()

	expected := fmt.Sprintf("%s:%d", defaultSMTPHost, utils.MockSMTPPort)
	SetSMTPPort(utils.MockSMTPPort)
	smtpConf := SMTP()
	assert.Equal(t, expected, smtpConf.Addr())

	err = smtpConf.DialAndQuit()
	assert.NoError(t, err)
}

func TestSMTPWithoutSMTPServer(t *testing.T) {
	SetSMTPPort(utils.MockSMTPPort)
	smtpConf := SMTP()

	err := smtpConf.DialAndQuit()
	expectedError := fmt.Sprintf(
		"smtp server dial: dial tcp %s: connect: connection refused",
		smtpConf.Addr())
	assert.EqualError(t, err, expectedError)
}
