package handler

import (
	"os"
	"testing"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	os.Setenv("AUTH_TEST", "true")
	dbConf := configs.DB()
	db.ResetDB(dbConf.TCPConnectionString(), dbConf.DBNameForTest())
	db.Sync(dbConf.ConnectionString(), dbConf.Echo)
}

func teardown() {
}
