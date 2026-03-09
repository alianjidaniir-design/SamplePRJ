package mysql

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultTaskTableName          = "tasks"
	defaultMaxOpenConnections     = 10
	defaultMaxIdleConnections     = 5
	defaultConnMaxLifetimeSeconds = 300
)

type Config struct {
	DSN                    string
	TaskTableName          string
	MaxOpenConnections     int
	MaxIdleConnections     int
	ConnMaxLifetimeSeconds int
}

func LoadConfigFromEnv() (cfg Config, err error) {
	cfg = Config{
		DSN:                    normalizeDSN(strings.TrimSpace(os.Getenv("MYSQL_DSN"))),
		TaskTableName:          strings.TrimSpace(os.Getenv("MYSQL_TASK_TABLE")),
		MaxOpenConnections:     readEnvInt("MYSQL_MAX_OPEN_CONNS", defaultMaxOpenConnections),
		MaxIdleConnections:     readEnvInt("MYSQL_MAX_IDLE_CONNS", defaultMaxIdleConnections),
		ConnMaxLifetimeSeconds: readEnvInt("MYSQL_CONN_MAX_LIFETIME_SECONDS", defaultConnMaxLifetimeSeconds),
	}

	if cfg.TaskTableName == "" {
		cfg.TaskTableName = defaultTaskTableName
	}

	if err := ValidateTableName(cfg.TaskTableName); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func readEnvInt(envName string, defaultValue int) int {
	raw := strings.TrimSpace(os.Getenv(envName))
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}

	return value
}

func normalizeDSN(dsn string) string {
	if dsn == "" {
		return ""
	}

	if !strings.Contains(dsn, "?") {
		return dsn + "?parseTime=true&loc=UTC&charset=utf8mb4"
	}

	base, queryPart, _ := strings.Cut(dsn, "?")
	queryValues, err := url.ParseQuery(queryPart)
	if err != nil {
		return dsn
	}

	if queryValues.Get("parseTime") == "" {
		queryValues.Set("parseTime", "true")
	}
	if queryValues.Get("loc") == "" {
		queryValues.Set("loc", "UTC")
	}
	if queryValues.Get("charset") == "" {
		queryValues.Set("charset", "utf8mb4")
	}

	return fmt.Sprintf("%s?%s", base, queryValues.Encode())
}
