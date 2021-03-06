package main

import (
	"log"
	"os"
	"syscall"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
)

var testDBCon *gorm.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	configs.SetMode(configs.TestMode)
	dbConf, err := configs.DB()
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Reset(dbConf.DSN(), dbConf.DBName())
	if err != nil {
		log.Fatalln(err)
	}
	testDBCon, err = db.Connection(dbConf.DSN(), dbConf.Echo)
	if err != nil {
		log.Fatalln(err)
	}
}

func teardown() {
}

func TestFuncMain(t *testing.T) {
	go func() {
		main()
	}()

	conf := configs.App()
	for !isListen(localHost, conf.ListenPort) {
		continue
	}
	Quit <- syscall.SIGINT

	dbConf, err := configs.DB()
	assert.NoError(t, err)
	if dbConf.SyncModels {
		assert.True(t, testDBCon.HasTable("users"))
	} else {
		assert.False(t, testDBCon.HasTable("users"))
	}
}

func TestFuncMainWithDBSync(t *testing.T) {
	enKey := configs.EnvPrefix + "DB_SYNC_MODELS"
	os.Setenv(enKey, "1")
	go func() {
		main()
	}()

	conf := configs.App()
	for !isListen(localHost, conf.ListenPort) {
		continue
	}
	Quit <- syscall.SIGINT
	os.Unsetenv(enKey)
	assert.True(t, testDBCon.HasTable("users"))
}
