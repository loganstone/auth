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

// DB ...
func DB() *DatabaseConfigs {
	conf := DatabaseConfigs{
		host: defaultDBHost,
		port: defaultDBPort,
	}

	required := map[string]*string{
		envPrefix + "DB_ID":   &conf.id,
		envPrefix + "DB_PW":   &conf.pw,
		envPrefix + "DB_NAME": &conf.name,
	}

	notSet := make([]string, 0, len(required))
	for k, p := range required {
		notSet = append(notSet, k)
		if v, ok := os.LookupEnv(k); ok {
			trimmedV := strings.TrimSpace(v)
			if trimmedV != "" {
				*p = trimmedV
				notSet = notSet[:len(notSet)-1]
			}
		}
	}
	if len(notSet) > 0 {
		log.Fatalf(failedToLookup, strings.Join(notSet, ", "))
	}

	for k, p := range map[string]interface{}{
		envPrefix + "DB_HOST":      &conf.host,
		envPrefix + "DB_PORT":      &conf.port,
		envPrefix + "DB_ECHO":      &conf.Echo,
		envPrefix + "DB_AUTO_SYNC": &conf.AutoSync,
	} {
		if v, ok := os.LookupEnv(k); ok {
			switch pt := p.(type) {
			case *string:
				*pt = v
			case *bool:
				*pt = (v == "1" || strings.ToLower(v) == "true")
			default:
				log.Fatalf("unknow type %T\n", pt)
			}
		}
	}

	return &conf
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
	if c.siginupURL == "" {
		return ""
	}

	if c.siginupURL == defaultSignupURL {
		return fmt.Sprintf(c.siginupURL, c.ListenPort, token)
	}

	last := c.siginupURL[len(c.siginupURL)-1]
	if string(last) != "/" {
		token = "/" + token
	}
	return fmt.Sprintf("%s%s", c.siginupURL, token)
}

// SecretKeyLen .
func (c *AppConfigs) SecretKeyLen() int {
	return c.secretKeyLen
}

// GracefulShutdownTimeout .
func (c *AppConfigs) GracefulShutdownTimeout() time.Duration {
	return c.gracefulShutdownTimeout
}

// App .
func App() *AppConfigs {
	conf := AppConfigs{
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

	for k, p := range map[string]interface{}{
		envPrefix + "LISTEN_PORT":          &conf.ListenPort,
		envPrefix + "SIGNUP_TOKEN_EXPIRE":  &conf.SignupTokenExpire,
		envPrefix + "SESSION_TOKEN_EXPIRE": &conf.SessionTokenExpire,
		envPrefix + "JWT_SIGNIN_KEY":       &conf.JWTSigninKey,
		envPrefix + "ORG":                  &conf.Org,
		envPrefix + "SUPPORT_EMAIL":        &conf.SupportEmail,
		envPrefix + "PAGE_SIZE":            &conf.PageSize,
		envPrefix + "PAGE_SIZE_LIMIT":      &conf.PageSizeLimit,
		envPrefix + "SIGNUP_URL":           &conf.siginupURL,
	} {
		if v, ok := os.LookupEnv(k); ok {
			switch pt := p.(type) {
			case *string:
				*pt = v
			case *int:
				if i, err := strconv.Atoi(v); err == nil {
					*pt = i
				}
			default:
				log.Fatalf("unknow type %T\n", pt)
			}
		}
	}

	return &conf
}
