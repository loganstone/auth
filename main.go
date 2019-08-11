package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/router"
	"github.com/loganstone/auth/validator"
)

func main() {
	options := configs.Opts()

	db.Sync()

	// Echo instance
	e := echo.New()

	e.Validator = validator.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	router.Init(e)

	// Debug uri - /debug/pprof/
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	// Start server
	go func() {
		listen := fmt.Sprintf(":%d", options.PortToListen)
		if err := e.Start(listen); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(
		context.Background(),
		configs.TimeoutToGracefulShutdown*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
