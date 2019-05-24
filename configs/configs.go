package configs

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

// DefaultPort ...
const DefaultPort = 9900

// DatabaseConfigs ...
type DatabaseConfigs struct {
	ID   string
	PW   string
	Name string
}

// DB ...
func DB() *DatabaseConfigs {
	const errMsgFmt = "'%s' environment variable is required"

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
	return &DatabaseConfigs{id, pw, name}
}
