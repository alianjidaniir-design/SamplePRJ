package mySqlDS

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultsStudentTableName     = "studentss"
	defaultMaxOpenConnections    = 10
	defaultMaxIdleConnections    = 5
	defaultConnMaxLifetimeSecond = 300
)

type Config struct {
	DSN                 string
	StudentTableName    string
	MaxOpenConnections  int
	MaxIdleConnections  int
	ConnMaxLifetimeSpan int
}

func LoadConfigFromEnv() (cfg Config, err error) {
	cfg = Config{
		DSN:                 normalizeDSN(strings.TrimSpace(os.Getenv("MYSQL_DSN"))),
		StudentTableName:    strings.TrimSpace(os.Getenv("MYSQL_STUDENTS_TABLE")),
		MaxOpenConnections:  readEnvInt("MYSQL_MAX_OPEN_CONNECTIONS", defaultMaxOpenConnections),
		MaxIdleConnections:  readEnvInt("MYSQL_MAX_IDLE_CONNECTIONS", defaultMaxIdleConnections),
		ConnMaxLifetimeSpan: readEnvInt("MYSQL_CONN_MAX_LIFETIME_SPAN", defaultConnMaxLifetimeSecond),
	}
	if cfg.StudentTableName == "" {
		cfg.StudentTableName = defaultsStudentTableName
	}

	if err := ValidateTableName(cfg.StudentTableName); err != nil {
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
		return dsn + "?parseTime=true&loc=Asia%2FTehran&time_zone=%27%2B03:30%27&charset=utf8mb4"
	}
	base, queryParse, _ := strings.Cut(dsn, "?")
	queryValues, err := url.ParseQuery(queryParse)
	if err != nil {
		return dsn
	}

	if queryValues.Get("parseTime") == "" {
		queryValues.Set("parseTime", "true")
	}
	if queryValues.Get("loc") == "" {
		queryValues.Set("loc", "Asia/Tehran")
	}
	if queryValues.Get("time_zone") == "" {
		queryValues.Set("time_zone", "'+03:30'")
	}
	if queryValues.Get("charset") == "" {
		queryValues.Set("charset", "utf8mb4")
	}

	return fmt.Sprintf("%s?%s", base, queryValues.Encode())

}
