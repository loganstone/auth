package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/jinzhu/gorm/dialects/mysql" // driver
)

const (
	failedToLookup = "must set '%s' environment variable\n"
	dbConOpt       = "charset=utf8mb4&parseTime=True&loc=Local"
	dbConStr       = "%s:%s@/%s?%s"
	dbTCPConStr    = "%s:%s@tcp(%s:%s)/"
	envPrefix      = "AUTH_"
)

const (
	defaultListenPort         = 9999
	defaultSignupTokenExpire  = 1800 // 30 minutes
	defaultSessionTokenExpire = 3600 // 60 minutes
	defaultJWTSigninKey       = "PlzSetYourSigninKey"
	defaultOrg                = "Auth"
	defaultSupportEmail       = "auth@email.com"
	defaultPageSize           = "20"

	defaultSignupURL = "http://localhost:%d/signup/email/verification/%s"
)

const (
	defaultDBHost = "127.0.0.1"
	defaultDBPort = "3306"
)

// DatabaseConfigs ...
type DatabaseConfigs struct {
	id       string
	pw       string
	name     string
	host     string
	port     string
	Echo     bool
	AutoSync bool
}

// DBNameForTest .
func (c *DatabaseConfigs) DBNameForTest() string {
	return fmt.Sprintf("%s_test", c.name)
}

// ConnectionString .
func (c *DatabaseConfigs) ConnectionString() string {
	dbName := c.name
	if gin.Mode() == gin.TestMode {
		dbName = c.DBNameForTest()
	}
	return fmt.Sprintf(dbConStr, c.id, c.pw, dbName, dbConOpt)
}

// TCPConnectionString .
func (c *DatabaseConfigs) TCPConnectionString() string {
	return fmt.Sprintf(dbTCPConStr, c.id, c.pw, c.host, c.port)
}

var dbConfigs = DatabaseConfigs{
	host: defaultDBHost,
	port: defaultDBPort,
}

var requiredDBConf = map[string]*string{
	envPrefix + "DB_ID":   &dbConfigs.id,
	envPrefix + "DB_PW":   &dbConfigs.pw,
	envPrefix + "DB_NAME": &dbConfigs.name,
}

// DB ...
func DB() *DatabaseConfigs {
	for k, p := range requiredDBConf {
		v, ok := os.LookupEnv(k)
		if !ok {
			log.Fatalf(failedToLookup, k)
		}
		*p = v
	}

	if h, ok := os.LookupEnv(envPrefix + "DB_HOST"); ok {
		dbConfigs.host = h
	}

	if p, ok := os.LookupEnv(envPrefix + "DB_PORT"); ok {
		dbConfigs.port = p
	}

	if s, ok := os.LookupEnv(envPrefix + "DB_AUTO_SYNC"); ok {
		dbConfigs.AutoSync = (s == "1" || strings.ToLower(s) == "true")
	}

	echo := os.Getenv(envPrefix + "DB_ECHO")
	dbConfigs.Echo = (echo == "1" || strings.ToLower(echo) == "true")

	return &dbConfigs
}

// AppConfigs .
type AppConfigs struct {
	gracefulShutdownTimeout time.Duration

	ListenPort         int
	SignupTokenExpire  int
	SessionTokenExpire int
	JWTSigninKey       string
	Org                string
	SupportEmail       string
	PageSize           string
	PageSizeLimit      int

	secretKeyLen int

	siginupURL string
}

// SignupURL .
func (c *AppConfigs) SignupURL(token string) string {
	if appConfigs.siginupURL == "" {
		return ""
	}

	if appConfigs.siginupURL == defaultSignupURL {
		return fmt.Sprintf(appConfigs.siginupURL, appConfigs.ListenPort, token)
	}

	last := appConfigs.siginupURL[len(appConfigs.siginupURL)-1]
	if string(last) != "/" {
		token = "/" + token
	}
	return fmt.Sprintf("%s%s", appConfigs.siginupURL, token)
}

// SecretKeyLen .
func (c *AppConfigs) SecretKeyLen() int {
	return c.secretKeyLen
}

// GracefulShutdownTimeout .
func (c *AppConfigs) GracefulShutdownTimeout() time.Duration {
	return c.gracefulShutdownTimeout
}

var appConfigs = AppConfigs{
	gracefulShutdownTimeout: 5,
	secretKeyLen:            16,

	ListenPort:         defaultListenPort,
	SignupTokenExpire:  defaultSignupTokenExpire,
	SessionTokenExpire: defaultSessionTokenExpire,
	JWTSigninKey:       defaultJWTSigninKey,
	Org:                defaultOrg,
	SupportEmail:       defaultSupportEmail,
	PageSize:           defaultPageSize,
	siginupURL:         defaultSignupURL,
}

var stringAppConf = map[string]*string{
	envPrefix + "JWT_SIGNIN_KEY": &appConfigs.JWTSigninKey,
	envPrefix + "ORG":            &appConfigs.Org,
	envPrefix + "SUPPORT_EMAIL":  &appConfigs.SupportEmail,
	envPrefix + "SIGNUP_URL":     &appConfigs.siginupURL,
	envPrefix + "PAGE_SIZE":      &appConfigs.PageSize,
}

var intAppConf = map[string]*int{
	envPrefix + "LISTEN_PORT":          &appConfigs.ListenPort,
	envPrefix + "SIGNUP_TOKEN_EXPIRE":  &appConfigs.SignupTokenExpire,
	envPrefix + "SESSION_TOKEN_EXPIRE": &appConfigs.SessionTokenExpire,
	envPrefix + "PAGE_SIZE_LIMIT":      &appConfigs.PageSizeLimit,
}

// App .
func App() *AppConfigs {
	for k, p := range stringAppConf {
		if v, ok := os.LookupEnv(k); ok {
			*p = v
		}
	}

	for k, p := range intAppConf {
		if v, ok := os.LookupEnv(k); ok {
			if i, err := strconv.Atoi(v); err == nil {
				*p = i
			}
		}
	}

	return &appConfigs
}
