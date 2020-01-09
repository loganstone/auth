package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/handler"
)

const localHost = "localhost"

// Quit .
var Quit = make(chan os.Signal)

func isListen(host string, port int) bool {
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

func main() {
	dbConf, err := configs.DB()
	if err != nil {
		log.Fatalln(err)
	}
	if dbConf.AutoSync {
		log.Println("Sync DB ...")
		db.Sync(dbConf.ConnectionString(), dbConf.Echo)
		log.Println("Sync DB Completed")
	}

	conf := configs.App()
	if isListen(localHost, conf.ListenPort) {
		log.Fatalf(`'%d' port already in use
- using env: export %sLISTEN_PORT=<port not in use>
`, conf.ListenPort, configs.EnvPrefix)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.ListenPort),
		Handler: handler.New(),
	}

	go func() {
		log.Printf("listen port: %d\n", conf.ListenPort)
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		conf.GracefulShutdownTimeout()*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
