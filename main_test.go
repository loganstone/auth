package main

import (
	"net"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

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
	// TODO(hs.lee):
	// 테스트 시 AUTH_DB_NAME 가 유지 되도록 처리
	os.Setenv("AUTH_TEST", "true")
	dbConf := configs.DB()
	db.ResetTestDB(dbConf.TCPConnectionString())
	DBSync()
}

func teardown() {
}

func TestFuncMain(t *testing.T) {
	go func() {
		main()
	}()

	conf := configs.App()
	for !IsListen("localhost", conf.PortToListen) {
		continue
	}
	Quit <- syscall.SIGINT
}

func IsListen(host string, port int) bool {
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(host, strconv.Itoa(port)),
		time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
