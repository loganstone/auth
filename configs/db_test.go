package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	SetMode(TestMode)
}

func TestDB(t *testing.T) {
	for k, v := range map[string]string{
		EnvPrefix + "DB_ID":          "test_db_id",
		EnvPrefix + "DB_PW":          "test_db_pw",
		EnvPrefix + "DB_NAME":        "test_db_name",
		EnvPrefix + "DB_HOST":        "127.0.0.1",
		EnvPrefix + "DB_PORT":        "3306",
		EnvPrefix + "DB_ECHO":        "false",
		EnvPrefix + "DB_SYNC_MODELS": "false",
	} {
		os.Setenv(k, v)
	}

	conf, err := DB()
	assert.Nil(t, err)

	expected := "test_db_name_test"
	dbName := conf.DBName()
	assert.Equal(t, expected, dbName)

	expected = "test_db_id:test_db_pw@(127.0.0.1:3306)/test_db_name_test?" + dbConOpt
	assert.Equal(t, expected, conf.ConnectionString())

	expected = "test_db_id:test_db_pw@tcp(127.0.0.1:3306)/"
	assert.Equal(t, expected, conf.TCPConnectionString())

	assert.False(t, conf.Echo)
	assert.False(t, conf.SyncModels)
}

func TestDBWithDebugMode(t *testing.T) {
	SetMode(DebugMode)

	for k, v := range map[string]string{
		EnvPrefix + "DB_ID":          "test_db_id",
		EnvPrefix + "DB_PW":          "test_db_pw",
		EnvPrefix + "DB_NAME":        "test_db_name",
		EnvPrefix + "DB_HOST":        "127.0.0.1",
		EnvPrefix + "DB_PORT":        "3306",
		EnvPrefix + "DB_ECHO":        "false",
		EnvPrefix + "DB_SYNC_MODELS": "false",
	} {
		os.Setenv(k, v)
	}

	conf, err := DB()
	assert.Nil(t, err)

	expected := "test_db_name"
	dbName := conf.DBName()
	assert.Equal(t, expected, dbName)

	expected = "test_db_id:test_db_pw@(127.0.0.1:3306)/test_db_name?" + dbConOpt
	assert.Equal(t, expected, conf.ConnectionString())

	expected = "test_db_id:test_db_pw@tcp(127.0.0.1:3306)/"
	assert.Equal(t, expected, conf.TCPConnectionString())

	assert.False(t, conf.Echo)
	assert.False(t, conf.SyncModels)

	SetMode(TestMode)
}

func TestDBWithMissingRequirement(t *testing.T) {
	missed := []string{
		EnvPrefix + "DB_ID",
		EnvPrefix + "DB_PW",
		EnvPrefix + "DB_NAME",
	}

	for _, v := range missed {
		os.Setenv(v, "")
	}

	conf, err := DB()
	assert.Nil(t, conf)

	expected := missingRequirementError("DB", missed)
	assert.Equal(t, expected, err)
}
