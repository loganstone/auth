package configs

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppDefault(t *testing.T) {
	conf := App()
	table := []struct {
		EnvName  string
		Expected interface{}
		Real     interface{}
	}{
		{
			EnvPrefix + "LISTEN_PORT",
			defaultListenPort,
			conf.ListenPort,
		},
		{
			EnvPrefix + "SIGNUP_TOKEN_EXPIRE",
			defaultSignupTokenExpire,
			conf.SignupTokenExpire,
		},
		{
			EnvPrefix + "SESSION_TOKEN_EXPIRE",
			defaultSessionTokenExpire,
			conf.SessionTokenExpire,
		},
		{
			EnvPrefix + "JWT_SIGNIN_KEY",
			defaultJWTSigninKey,
			conf.JWTSigninKey,
		},
		{
			EnvPrefix + "ORG",
			defaultOrg,
			conf.Org,
		},
		{
			EnvPrefix + "SUPPORT_EMAIL",
			defaultSupportEmail,
			conf.SupportEmail,
		},
		{
			EnvPrefix + "PAGE_SIZE",
			defaultPageSize,
			conf.PageSize,
		},
		{
			EnvPrefix + "PAGE_SIZE_LIMIT",
			0,
			conf.PageSizeLimit,
		},
	}

	for _, v := range table {
		if _, ok := os.LookupEnv(v.EnvName); !ok {
			assert.Equal(t, v.Expected, v.Real)
		}
	}

	assert.Equal(t, 16, conf.SecretKeyLen())
	assert.Equal(t, time.Duration(5), conf.GracefulShutdownTimeout())
}

func TestApp(t *testing.T) {
	data := map[string]string{
		EnvPrefix + "LISTEN_PORT":          "8080",
		EnvPrefix + "SIGNUP_TOKEN_EXPIRE":  "3600",
		EnvPrefix + "SESSION_TOKEN_EXPIRE": "3600",
		EnvPrefix + "JWT_SIGNIN_KEY":       "testkey",
		EnvPrefix + "ORG":                  "test org",
		EnvPrefix + "SUPPORT_EMAIL":        "test.support@email.com",
		EnvPrefix + "PAGE_SIZE":            "50",
		EnvPrefix + "PAGE_SIZE_LIMIT":      "100",
	}

	for k, v := range data {
		os.Setenv(k, v)
	}

	conf := App()
	val, err := strconv.Atoi(data[EnvPrefix+"LISTEN_PORT"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.ListenPort)

	assert.Equal(t, data[EnvPrefix+"JWT_SIGNIN_KEY"], conf.JWTSigninKey)

	val, err = strconv.Atoi(data[EnvPrefix+"SIGNUP_TOKEN_EXPIRE"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.SignupTokenExpire)

	val, err = strconv.Atoi(data[EnvPrefix+"SESSION_TOKEN_EXPIRE"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.SessionTokenExpire)

	assert.Equal(t, data[EnvPrefix+"ORG"], conf.Org)

	assert.Equal(t, data[EnvPrefix+"SUPPORT_EMAIL"], conf.SupportEmail)

	assert.Equal(t, data[EnvPrefix+"PAGE_SIZE"], conf.PageSize)

	val, err = strconv.Atoi(data[EnvPrefix+"PAGE_SIZE_LIMIT"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.PageSizeLimit)
}

func TestSignupURL(t *testing.T) {
	conf := App()
	token := "testtoken"
	expected := fmt.Sprintf(defaultSignupURL, conf.ListenPort, token)
	url := conf.SignupURL(token)
	assert.Equal(t, expected, url)
}

func TestSignupURLWithSetEnv(t *testing.T) {
	token := "testtoken"
	table := []struct {
		URL      string
		Expected string
	}{
		{"", ""},
		{"http://example.com/", "http://example.com/" + token},
		{"http://example.com", "http://example.com/" + token},
	}

	for _, v := range table {
		os.Setenv(EnvPrefix+"SIGNUP_URL", v.URL)
		conf := App()
		url := conf.SignupURL(token)
		assert.Equal(t, v.Expected, url)
	}
}
