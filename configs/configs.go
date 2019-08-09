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
	TimeoutToGracefulShutdown = 10

	connOpt     = "charset=utf8mb4&parseTime=True&loc=Local"
	defaultPort = 9900
	envErrFmt   = "'%s' environment variable is required\n"
)

// Options .
type Options struct {
	PortToListen int
}

var options Options

// DatabaseConfigs ...
type DatabaseConfigs struct {
	ID   string
	PW   string
	Name string
	Echo bool
}

// ConnectionString .
func (c *DatabaseConfigs) ConnectionString() string {
	confSlice := append([]interface{}{c.ID, c.PW, c.Name}, connOpt)
	return fmt.Sprintf("%s:%s@/%s?%s", confSlice...)
}

// DB ...
func DB() *DatabaseConfigs {
	id, ok := os.LookupEnv("AUTH_DB_ID")
	if !ok {
		log.Fatalf(envErrFmt, "AUTH_DB_ID")
	}

	pw, ok := os.LookupEnv("AUTH_DB_PW")
	if !ok {
		log.Fatalf(envErrFmt, "AUTH_DB_PW")
	}

	name, ok := os.LookupEnv("AUTH_DB_NAME")
	if !ok {
		log.Fatalf(envErrFmt, "AUTH_DB_NAME")
	}

	echo := os.Getenv("AUTH_DB_ECHO")
	return &DatabaseConfigs{
		id, pw, name, (echo == "1" || strings.ToLower(echo) == "true"),
	}
}

func init() {
	flag.IntVar(&options.PortToListen, "p", defaultPort, "port to listen on")
	flag.Parse()
}

// Opts .
func Opts() Options {
	return options
}
