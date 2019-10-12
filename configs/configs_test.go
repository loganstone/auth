package configs

import (
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
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
}

func TestApp(t *testing.T) {
	conf := App()
	assert.Equal(t, defaultPortToListen, conf.PortToListen)
	assert.Equal(t, defaultSignupTokenExpire, conf.SignupTokenExpire)
	assert.Equal(t, defaultSessionTokenExpire, conf.SessionTokenExpire)
	assert.Equal(t, []byte(defaultJWTSigninKey), conf.JWTSigninKey)
	assert.Equal(t, defaultPageSize, conf.PageSize)
}
