package main

import (
	"os"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/stretchr/testify/assert"
)

var testDBCon *gorm.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	gin.SetMode(gin.TestMode)
	dbConf := configs.DB()
	db.ResetDB(dbConf.TCPConnectionString(), dbConf.DBNameForTest())
	testDBCon = db.Connection(dbConf.ConnectionString(), dbConf.Echo)
}

func teardown() {
}

func TestFuncMain(t *testing.T) {
	go func() {
		main()
	}()

	conf := configs.App()
	for !isListen(localHost, conf.PortToListen) {
		continue
	}
	Quit <- syscall.SIGINT
	assert.False(t, testDBCon.HasTable("users"))
}

func TestFuncMainWithDBSync(t *testing.T) {
	os.Setenv("AUTH_DB_SYNC", "1")
	go func() {
		main()
	}()

	conf := configs.App()
	for !isListen(localHost, conf.PortToListen) {
		continue
	}
	Quit <- syscall.SIGINT
	os.Unsetenv("AUTH_DB_SYNC")
	assert.True(t, testDBCon.HasTable("users"))
}
