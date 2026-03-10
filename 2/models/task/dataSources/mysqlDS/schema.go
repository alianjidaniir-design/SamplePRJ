package mysqlDS

import (
	"database/sql"
	"fmt"
	"regexp"
)

var safeTableNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func ValidateTableName(tableName string) error {
	if !safeTableNamePattern.MatchString(tableName) {
		return fmt.Errorf("invalid table name %q: only letters, numbers, and underscore are allowed", tableName)
	}

	return nil
}

func TaskTableIdentifier(tableName string) (string, error) {
	if err := ValidateTableName(tableName); err != nil {
		return "", err
	}

	return fmt.Sprintf("`%s`", tableName), nil
}

func EnsureTaskTable(db *sql.DB, tableName string) error {
	tableIdentifier, err := TaskTableIdentifier(tableName)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id BIGINT NOT NULL AUTO_INCREMENT,
	title VARCHAR(128) NOT NULL,
	description VARCHAR(512) NOT NULL DEFAULT '',
	createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	INDEX idx_createdAt (createdAt)
);`, tableIdentifier)

	_, err = db.Exec(query)
	return err
}
