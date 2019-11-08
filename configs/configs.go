package configs

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql" // driver
)

const (
	// TimeoutToGracefulShutdown .
	TimeoutToGracefulShutdown = 5

	failedToLookup            = "must set '%s' environment variable\n"
	dbConOpt                  = "charset=utf8mb4&parseTime=True&loc=Local"
	defaultPortToListen       = 9999
	defaultSignupTokenExpire  = 1800 // 30 minutes
	defaultSessionTokenExpire = 3600 // 60 minutes
	defaultJWTSigninKey       = "PlzSetYourSigninKey"
	defaultPageSize           = "20"
	envPrefix                 = "AUTH_"
)

// AppConfigs .
type AppConfigs struct {
	PortToListen       int
	SignupTokenExpire  int
	SessionTokenExpire int
	JWTSigninKey       string
	PageSize           string
	PageSizeLimit      int
}

var appConfigs = AppConfigs{
	PortToListen:       defaultPortToListen,
	SignupTokenExpire:  defaultSignupTokenExpire,
	SessionTokenExpire: defaultSessionTokenExpire,
	JWTSigninKey:       defaultJWTSigninKey,
	PageSize:           defaultPageSize,
}

// DatabaseConfigs ...
type DatabaseConfigs struct {
	id   string
	pw   string
	name string
	Echo bool
}

// ConnectionString .
func (c *DatabaseConfigs) ConnectionString() string {
	return fmt.Sprintf("%s:%s@/%s?%s", c.id, c.pw, c.name, dbConOpt)
}

func init() {
	flag.IntVar(
		&appConfigs.PortToListen, "p",
		defaultPortToListen, "port to listen on")
}

// DB ...
func DB() *DatabaseConfigs {
	id, ok := os.LookupEnv(envPrefix + "DB_ID")
	if !ok {
		log.Fatalf(failedToLookup, envPrefix+"DB_ID")
	}

	pw, ok := os.LookupEnv(envPrefix + "DB_PW")
	if !ok {
		log.Fatalf(failedToLookup, envPrefix+"DB_PW")
	}

	name, ok := os.LookupEnv(envPrefix + "DB_NAME")
	if !ok {
		log.Fatalf(failedToLookup, envPrefix+"DB_NAME")
	}

	echo := os.Getenv(envPrefix + "DB_ECHO")
	return &DatabaseConfigs{
		id, pw, name, (echo == "1" || strings.ToLower(echo) == "true"),
	}
}

// App .
func App() *AppConfigs {
	if key, ok := os.LookupEnv(envPrefix + "JWT_SIGNIN_KEY"); ok {
		appConfigs.JWTSigninKey = key
	}

	if expire, ok := os.LookupEnv(envPrefix + "SIGNUP_TOKEN_EXPIRE"); ok {
		v, err := strconv.Atoi(expire)
		if err == nil {
			appConfigs.SignupTokenExpire = v
		}
	}

	if expire, ok := os.LookupEnv(envPrefix + "SESSION_TOKEN_EXPIRE"); ok {
		v, err := strconv.Atoi(expire)
		if err == nil {
			appConfigs.SessionTokenExpire = v
		}
	}

	if pageSize, ok := os.LookupEnv(envPrefix + "PAGE_SIZE"); ok {
		appConfigs.PageSize = pageSize
	}

	if pageSizeLimit, ok := os.LookupEnv(envPrefix + "PAGE_SIZE_LIMIT"); ok {
		if v, err := strconv.Atoi(pageSizeLimit); err == nil {
			appConfigs.PageSizeLimit = v
		}
	}

	return &appConfigs
}
