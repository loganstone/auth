package configs

import (
	"fmt"
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
	confSlice := append(conf.ToSlice(), ConnOpt)
	connectionString := fmt.Sprintf("%s:%s@/%s?%s", confSlice...)
	expected := "test_db_id:test_db_pw@/test_db_name?" + ConnOpt
	assert.Equal(t, connectionString, expected)
}
