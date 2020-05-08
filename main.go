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

var conf = configs.App()

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

func syncModels(c *configs.DatabaseConfig) error {
	log.Println("sync models start ...")
	con, err := db.SyncModels(c.DSN(), c.Echo)
	defer con.Close()
	if err != nil {
		return err
	}
	log.Println("sync models completed")
	return nil
}

func server() *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.ListenPort),
		Handler: handler.New(),
	}
}

func checkListenPort() {
	if isListen(localHost, conf.ListenPort) {
		log.Fatalf(`'%d' port already in use
- using env: export %sLISTEN_PORT=<port not in use>
`, conf.ListenPort, configs.EnvPrefix)
	}
}

func main() {
	if configs.Mode() != configs.TestMode {
		smtpConf := configs.SMTP()
		err := smtpConf.DialAndQuit()
		if err != nil {
			log.Fatalln(err)
		}
	}

	dbConf, err := configs.DB()
	if err != nil {
		log.Fatalln(err)
	}

	if dbConf.SyncModels {
		if err := syncModels(dbConf); err != nil {
			log.Fatalln(err)
		}
	}

	checkListenPort()

	srv := server()
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
	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		conf.GracefulShutdownDuration())
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}

	log.Println("server exiting")
}
