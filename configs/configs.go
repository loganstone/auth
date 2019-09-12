package configs

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

const (
	// TimeoutToGracefulShutdown .
	TimeoutToGracefulShutdown = 5

	dbConOpt            = "charset=utf8mb4&parseTime=True&loc=Local"
	defaultPortToListen = 9900
	failedToLookup      = "need to set '%s' environment variable\n"
)

// AppConfigs .
type AppConfigs struct {
	PortToListen int
}

var appConfigs AppConfigs

// DatabaseConfigs ...
type DatabaseConfigs struct {
	id   string
	pw   string
	name string
	Echo bool
}

// ConnectionString .
func (c *DatabaseConfigs) ConnectionString() string {
	conf := append([]interface{}{c.id, c.pw, c.name}, dbConOpt)
	return fmt.Sprintf("%s:%s@/%s?%s", conf...)
}

func init() {
	flag.IntVar(
		&appConfigs.PortToListen, "p",
		defaultPortToListen, "port to listen on")
}

// DB ...
func DB() *DatabaseConfigs {
	id, ok := os.LookupEnv("AUTH_DB_ID")
	if !ok {
		log.Fatalf(failedToLookup, "AUTH_DB_ID")
	}

	pw, ok := os.LookupEnv("AUTH_DB_PW")
	if !ok {
		log.Fatalf(failedToLookup, "AUTH_DB_PW")
	}

	name, ok := os.LookupEnv("AUTH_DB_NAME")
	if !ok {
		log.Fatalf(failedToLookup, "AUTH_DB_NAME")
	}

	echo := os.Getenv("AUTH_DB_ECHO")
	return &DatabaseConfigs{
		id, pw, name, (echo == "1" || strings.ToLower(echo) == "true"),
	}
}

// App .
func App() AppConfigs {
	return appConfigs
}
