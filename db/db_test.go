package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/configs"
)

func TestConnection(t *testing.T) {
	configs.SetMode(configs.TestMode)
	dbConf, err := configs.DB()
	assert.NoError(t, err)
	con, err := Connection(dbConf.DSN(), true)
	assert.NotNil(t, con)
}

func TestConnectionWithBadDSN(t *testing.T) {
	baddsn := "baddsn"
	_, err := Connection(baddsn, true)
	expectedError := "invalid DSN: missing the slash separating the database name"
	assert.EqualError(t, err, expectedError)
}

func TestSyncModels(t *testing.T) {
	configs.SetMode(configs.TestMode)
	dbConf, err := configs.DB()
	assert.NoError(t, err)
	con, err := SyncModels(dbConf.DSN(), true)
	assert.NotNil(t, con)
}

func TestSyncModelsWithBadDSN(t *testing.T) {
	baddsn := "baddsn"
	_, err := SyncModels(baddsn, true)
	expectedError := "invalid DSN: missing the slash separating the database name"
	assert.EqualError(t, err, expectedError)
}

func TestReset(t *testing.T) {
	configs.SetMode(configs.TestMode)
	dbConf, err := configs.DB()
	err = Reset(dbConf.DSN(), dbConf.DBName())
	assert.NoError(t, err)
}

func TestResetWithBadDSN(t *testing.T) {
	err := Reset("baddsn", "badtable")
	expectedError := "db connection failed: invalid DSN: missing the slash separating the database name"
	assert.EqualError(t, err, expectedError)
}
