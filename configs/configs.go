package configs

import (
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
	defaultDBHost             = "127.0.0.1"
	defaultDBPort             = "3306"
	dbConStr                  = "%s:%s@/%s?%s"
	dbTCPConStr               = "%s:%s@tcp(%s:%s)/"
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
	Name string
	host string
	port string
	Echo bool
}

var dbConfigs = DatabaseConfigs{
	host: defaultDBHost,
	port: defaultDBPort,
}

// ConnectionString .
func (c *DatabaseConfigs) ConnectionString() string {
	// TODO(hs.lee):
	// main_test.go 에서 설정한
	// AUTH_TEST 환경변수를 모든 테스트에
	// 유지시킬 방법을 찾아봐야 한다.
	if v, ok := os.LookupEnv("AUTH_TEST"); ok {
		if v == "true" {
			return fmt.Sprintf(
				dbConStr, c.id, c.pw, c.Name+"_test", dbConOpt)
		}
	}
	return fmt.Sprintf(dbConStr, c.id, c.pw, c.Name, dbConOpt)
}

// TCPConnectionString .
func (c *DatabaseConfigs) TCPConnectionString() string {
	return fmt.Sprintf(dbTCPConStr, c.id, c.pw, c.host, c.port)
}

// DB ...
func DB() *DatabaseConfigs {
	if dbConfigs.id == "" {
		id, ok := os.LookupEnv(envPrefix + "DB_ID")
		if !ok {
			log.Fatalf(failedToLookup, envPrefix+"DB_ID")
		}
		dbConfigs.id = id
	}

	if dbConfigs.pw == "" {
		pw, ok := os.LookupEnv(envPrefix + "DB_PW")
		if !ok {
			log.Fatalf(failedToLookup, envPrefix+"DB_PW")
		}
		dbConfigs.pw = pw
	}

	if dbConfigs.Name == "" {
		name, ok := os.LookupEnv(envPrefix + "DB_NAME")
		if !ok {
			log.Fatalf(failedToLookup, envPrefix+"DB_NAME")
		}
		dbConfigs.Name = name
	}

	if h, ok := os.LookupEnv(envPrefix + "DB_HOST"); ok {
		dbConfigs.host = h
	}

	if p, ok := os.LookupEnv(envPrefix + "DB_PORT"); ok {
		dbConfigs.port = p
	}

	echo := os.Getenv(envPrefix + "DB_ECHO")
	dbConfigs.Echo = (echo == "1" || strings.ToLower(echo) == "true")

	return &dbConfigs
}

// App .
func App() *AppConfigs {
	if port, ok := os.LookupEnv(envPrefix + "LISTEN_PORT"); ok {
		v, err := strconv.Atoi(port)
		if err == nil {
			appConfigs.PortToListen = v
		}
	}

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
