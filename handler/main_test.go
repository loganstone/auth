package handler

import (
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"

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
	err = db.Reset(dbConf.TCPConnectionString(), dbConf.DBName())
	if err != nil {
		log.Fatalln(err)
	}
	testDBCon, err = db.SyncModels(dbConf.ConnectionString(), dbConf.Echo)
	if err != nil {
		log.Fatalln(err)
	}
}

func teardown() {
	testDBCon.Close()
}
