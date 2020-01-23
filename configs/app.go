package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
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