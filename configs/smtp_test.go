package configs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTP(t *testing.T) {
	expected := fmt.Sprintf("%s:%s", defaultSMTPHost, defaultSMTPPort)
	smtpConf := SMTP()
	assert.Equal(t, expected, smtpConf.Addr())

	err := smtpConf.DialAndQuit()
	assert.NoError(t, err)
}
