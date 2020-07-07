package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultGracefulShutdownDuration = 5
	defaultSecretKeyLen             = 16

	defaultListenPort         = 9999
	defaultSignupTokenExpire  = 1800 // 30 minutes
	defaultSessionTokenExpire = 3600 // 60 minutes
	defaultJWTSigninKey       = "PlzSetYourSigninKey"
	defaultOrg                = "Auth"
	defaultSupportEmail       = "auth@email.com"
	defaultPageSize           = "20"

	defaultSignupURL = "http://localhost:%d/signup/email/verification/%s"
)

// AppConfig contains the values needed to operate application.
type AppConfig struct {
	gracefulShutdownDuration time.Duration

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

// SignupURL is returns signup url to be used by frontend.
func (c *AppConfig) SignupURL(token string) string {
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

// SecretKeyLen is returns key length value required when creating a secretKey.
func (c *AppConfig) SecretKeyLen() int {
	return c.secretKeyLen
}

// GracefulShutdownDuration is returns duration value to be used when GracefulShutdown.
func (c *AppConfig) GracefulShutdownDuration() time.Duration {
	return c.gracefulShutdownDuration
}

// App .
func App() *AppConfig {
	conf := AppConfig{
		gracefulShutdownDuration: time.Second * defaultGracefulShutdownDuration,
		secretKeyLen:             defaultSecretKeyLen,

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
