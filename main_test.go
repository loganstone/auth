package main

import (
	"net"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/loganstone/auth/configs"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// TODO(hs.lee):
	// rest test database 추가
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
