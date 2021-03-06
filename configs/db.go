package configs

import (
	"fmt"
	"os"
	"strings"
)

const (
	dbConOpt = "charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	dbConStr = "%s:%s@(%s:%s)/%s?%s"
)

const (
	defaultDBHost = "127.0.0.1"
	defaultDBPort = "3306"
)

// DatabaseConfig contains values for database access.
type DatabaseConfig struct {
	id         string
	pw         string
	name       string
	host       string
	port       string
	Echo       bool
	SyncModels bool
}

// DBName is returns database name to be used in the current application.
func (c *DatabaseConfig) DBName() string {
	if Mode() == TestMode {
		return fmt.Sprintf("%s_test", c.name)
	}
	return c.name
}

// DSN is returns database source name.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		dbConStr, c.id, c.pw, c.host, c.port, c.DBName(), dbConOpt)
}

// DB returns the values needed to access the database.
// If required value constraint is not met, an error is returned.
func DB() (*DatabaseConfig, error) {
	const fnDB = "DB"
	conf := DatabaseConfig{
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
		EnvPrefix + "DB_HOST":        &conf.host,
		EnvPrefix + "DB_PORT":        &conf.port,
		EnvPrefix + "DB_ECHO":        &conf.Echo,
		EnvPrefix + "DB_SYNC_MODELS": &conf.SyncModels,
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
