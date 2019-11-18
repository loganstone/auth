package main

import (
	"os"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"

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
	gin.SetMode(gin.TestMode)
	dbConf := configs.DB()
	db.ResetDB(dbConf.TCPConnectionString(), dbConf.DBNameForTest())
	db.Sync(dbConf.ConnectionString(), dbConf.Echo)
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
}
