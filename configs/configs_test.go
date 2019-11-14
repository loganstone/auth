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

	expected = "test_db_id:test_db_pw@tcp(127.0.0.1:3306)/"
	assert.Equal(t, conf.TCPConnectionString(), expected)
}

func TestApp(t *testing.T) {
	// TODO(hs.lee):
	// 올바른 테스트가 아니다.
	// 다시 작성 할것
	conf := App()
	assert.Equal(t, defaultSignupTokenExpire, conf.SignupTokenExpire)
	assert.Equal(t, defaultSessionTokenExpire, conf.SessionTokenExpire)
	assert.Equal(t, defaultJWTSigninKey, conf.JWTSigninKey)
	assert.Equal(t, defaultPageSize, conf.PageSize)
}
