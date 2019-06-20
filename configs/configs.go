package configs

import (
	"log"
	"os"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

// ConnOpt .
const ConnOpt = "charset=utf8mb4&parseTime=True&loc=Local"

// DefaultPort ...
const DefaultPort = 9900

// TimeoutToGracefulShutdown ...
const TimeoutToGracefulShutdown = 10

// DatabaseConfigs ...
type DatabaseConfigs struct {
	ID   string
	PW   string
	Name string
	Echo bool
}

// ToSlice .
func (c *DatabaseConfigs) ToSlice() []interface{} {
	return []interface{}{c.ID, c.PW, c.Name}
}

// DB ...
func DB() *DatabaseConfigs {
	const errMsgFmt = "'%s' environment variable is required\n"

	id, ok := os.LookupEnv("AUTH_DB_ID")
	if !ok {
		log.Fatalf(errMsgFmt, "AUTH_DB_ID")
	}

	pw, ok := os.LookupEnv("AUTH_DB_PW")
	if !ok {
		log.Fatalf(errMsgFmt, "AUTH_DB_PW")
	}

	name, ok := os.LookupEnv("AUTH_DB_NAME")
	if !ok {
		log.Fatalf(errMsgFmt, "AUTH_DB_NAME")
	}

	echo := os.Getenv("AUTH_DB_ECHO")
	return &DatabaseConfigs{
		id, pw, name, (echo == "1" || strings.ToLower(echo) == "true"),
	}
}
