package handler

import (
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
	dbConf := configs.DB()
	db.ResetDB(dbConf.TCPConnectionString(), dbConf.DBNameForTest())
	db.Sync(dbConf.ConnectionString(), dbConf.Echo)
	testDBCon = GetDBConnection()
}

func teardown() {
	testDBCon.Close()
}
