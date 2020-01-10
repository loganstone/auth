package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/jinzhu/gorm/dialects/mysql" // driver
)

const (
	// EnvPrefix .
	EnvPrefix      = "AUTH_"
	failedToLookup = "must set '%s' environment variable"
	dbConOpt       = "charset=utf8mb4&parseTime=True&loc=Local"
	dbConStr       = "%s:%s@/%s?%s"
	dbTCPConStr    = "%s:%s@tcp(%s:%s)/"
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

// EnvError .
type EnvError struct {
	Func string
	Err  error
}

func (e *EnvError) Error() string {
	return "configs." + e.Func + ": " + e.Err.Error()
}

func missingRequirementError(fn string, missed []string) *EnvError {
	err := fmt.Sprintf(failedToLookup, strings.Join(missed, ", "))
	return &EnvError{fn, errors.New(err)}
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

// DB .
func DB() (*DatabaseConfigs, error) {
	const fnDB = "DB"
	conf := DatabaseConfigs{
		host: defaultDBHost,
		port: defaultDBPort,
	}

	required := []struct {
		EnvName string
		ConfRef *string
	}{
		{EnvPrefix + "DB_ID", &conf.id},
		{EnvPrefix + "DB_PW", &conf.pw},
		{EnvPrefix + "DB_NAME", &conf.name},
	}

	missed := make([]string, 0, len(required))
	for _, item := range required {
		missed = append(missed, item.EnvName)
		if v, ok := os.LookupEnv(item.EnvName); ok {
			v = strings.TrimSpace(v)
			if v != "" {
				*item.ConfRef = v
				missed = missed[:len(missed)-1]
			}
		}
	}
	if len(missed) > 0 {
		return nil, missingRequirementError(fnDB, missed)
	}

	for k, p := range map[string]interface{}{
		EnvPrefix + "DB_HOST":      &conf.host,
		EnvPrefix + "DB_PORT":      &conf.port,
		EnvPrefix + "DB_ECHO":      &conf.Echo,
		EnvPrefix + "DB_AUTO_SYNC": &conf.AutoSync,
	} {
		if v, ok := os.LookupEnv(k); ok {
			switch pt := p.(type) {
			case *string:
				*pt = v
			case *bool:
				*pt = (v == "1" || strings.ToLower(v) == "true")
			}
		}
	}

	return &conf, nil
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
		EnvPrefix + "LISTEN_PORT":          &conf.ListenPort,
		EnvPrefix + "SIGNUP_TOKEN_EXPIRE":  &conf.SignupTokenExpire,
		EnvPrefix + "SESSION_TOKEN_EXPIRE": &conf.SessionTokenExpire,
		EnvPrefix + "JWT_SIGNIN_KEY":       &conf.JWTSigninKey,
		EnvPrefix + "ORG":                  &conf.Org,
		EnvPrefix + "SUPPORT_EMAIL":        &conf.SupportEmail,
		EnvPrefix + "PAGE_SIZE":            &conf.PageSize,
		EnvPrefix + "PAGE_SIZE_LIMIT":      &conf.PageSizeLimit,
		EnvPrefix + "SIGNUP_URL":           &conf.siginupURL,
	} {
		if v, ok := os.LookupEnv(k); ok {
			switch pt := p.(type) {
			case *string:
				*pt = v
			case *int:
				if i, err := strconv.Atoi(v); err == nil {
					*pt = i
				}
			}
		}
	}

	return &conf
}
