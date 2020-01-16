package handler

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
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
	gin.SetMode(gin.TestMode)
	dbConf, err := configs.DB()
	if err != nil {
		log.Fatalln(err)
	}
	db.Reset(dbConf.TCPConnectionString(), dbConf.DBNameForTest())
	db.Sync(dbConf.ConnectionString(), dbConf.Echo)
	testDBCon = DBConnection()
}

func teardown() {
	testDBCon.Close()
}
