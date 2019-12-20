package configs

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	// Setup
	os.Setenv("AUTH_DB_ID", "test_db_id")
	os.Setenv("AUTH_DB_PW", "test_db_pw")
	os.Setenv("AUTH_DB_NAME", "test_db_name")
	os.Setenv("AUTH_DB_ECHO", "false")

	// Assertions
	conf := DB()
	expected := "test_db_id:test_db_pw@/test_db_name?" + dbConOpt
	assert.Equal(t, conf.ConnectionString(), expected)

	expected = "test_db_id:test_db_pw@tcp(127.0.0.1:3306)/"
	assert.Equal(t, conf.TCPConnectionString(), expected)
}

func TestAppDefault(t *testing.T) {
	conf := App()
	if _, ok := os.LookupEnv(envPrefix + "LISTEN_PORT"); !ok {
		assert.Equal(t, defaultListenPort, conf.ListenPort)
	}

	if _, ok := os.LookupEnv(envPrefix + "SIGNUP_TOKEN_EXPIRE"); !ok {
		assert.Equal(t, defaultSignupTokenExpire, conf.SignupTokenExpire)
	}

	if _, ok := os.LookupEnv(envPrefix + "SESSION_TOKEN_EXPIRE"); !ok {
		assert.Equal(t, defaultSessionTokenExpire, conf.SessionTokenExpire)
	}

	if _, ok := os.LookupEnv(envPrefix + "JWT_SIGNIN_KEY"); !ok {
		assert.Equal(t, defaultJWTSigninKey, conf.JWTSigninKey)
	}

	if _, ok := os.LookupEnv(envPrefix + "ORG"); !ok {
		assert.Equal(t, defaultOrg, conf.Org)
	}

	if _, ok := os.LookupEnv(envPrefix + "SUPPORT_EMAIL"); !ok {
		assert.Equal(t, defaultSupportEmail, conf.SupportEmail)
	}

	if _, ok := os.LookupEnv(envPrefix + "PAGE_SIZE"); !ok {
		assert.Equal(t, defaultPageSize, conf.PageSize)
	}

	if _, ok := os.LookupEnv(envPrefix + "PAGE_SIZE_LIMIT"); !ok {
		assert.Equal(t, 0, conf.PageSizeLimit)
	}
}

func TestApp(t *testing.T) {
	data := map[string]string{
		envPrefix + "LISTEN_PORT":          "8080",
		envPrefix + "JWT_SIGNIN_KEY":       "testkey",
		envPrefix + "SIGNUP_TOKEN_EXPIRE":  "3600",
		envPrefix + "SESSION_TOKEN_EXPIRE": "3600",
		envPrefix + "ORG":                  "test org",
		envPrefix + "SUPPORT_EMAIL":        "test.support@email.com",
		envPrefix + "PAGE_SIZE":            "50",
		envPrefix + "PAGE_SIZE_LIMIT":      "100",
	}

	for k, v := range data {
		os.Setenv(k, v)
	}

	conf := App()
	val, err := strconv.Atoi(data[envPrefix+"LISTEN_PORT"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.ListenPort)

	assert.Equal(t, data[envPrefix+"JWT_SIGNIN_KEY"], conf.JWTSigninKey)

	val, err = strconv.Atoi(data[envPrefix+"SIGNUP_TOKEN_EXPIRE"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.SignupTokenExpire)

	val, err = strconv.Atoi(data[envPrefix+"SESSION_TOKEN_EXPIRE"])
	assert.Nil(t, err)
	assert.Equal(t, val, conf.SessionTokenExpire)

	assert.Equal(t, data[envPrefix+"ORG"], conf.Org)

	assert.Equal(t, data[envPrefix+"SUPPORT_EMAIL"], conf.SupportEmail)

	assert.Equal(t, data[envPrefix+"PAGE_SIZE"], conf.PageSize)

	val, err = strconv.Atoi(data[envPrefix+"PAGE_SIZE_LIMIT"])
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
		os.Setenv(envPrefix+"SIGNUP_URL", v.URL)
		conf := App()
		url := conf.SignupURL(token)
		assert.Equal(t, v.Expected, url)
	}
}
