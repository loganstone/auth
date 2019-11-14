package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/handler"
)

// Quit .
var Quit = make(chan os.Signal)

func main() {
	// TODO(hs.lee):
	// 환경 변수로 설정 여부를 지정 하도록 변경.
	dbConf := configs.DB()
	db.Sync(dbConf.ConnectionString(), dbConf.Echo)

	conf := configs.App()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.PortToListen),
		Handler: handler.New(),
	}

	go func() {
		log.Printf("listen port: %d\n", conf.PortToListen)
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
		configs.TimeoutToGracefulShutdown*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Printf("timeout of %d seconds.", configs.TimeoutToGracefulShutdown)
	}
	log.Println("Server exiting")
}
