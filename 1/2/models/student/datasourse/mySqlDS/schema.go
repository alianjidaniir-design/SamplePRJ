package mySqlDS

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

func studentTableIdentifier(tableName string) (string, error) {
	if err := ValidateTableName(tableName); err != nil {
		return "", err

	}
	return fmt.Sprintf("`%s`", tableName), nil
}

func EnsureTaskTable(db *sql.DB, tableName string) error {
	tableIdentifier, err := studentTableIdentifier(tableName)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id BIGINT NOT NULL AUTO_INCREMENT,
    student_code VARCHAR(128) NOT NULL ,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(512) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL ,
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
);`, tableIdentifier)

	_, err = db.Exec(query)
	return err

}
