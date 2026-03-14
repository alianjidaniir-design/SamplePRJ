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
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NULL DEFAULT NULL,
	deleted_at TIMESTAMP NULL DEFAULT NULL,
	PRIMARY KEY (id),
	INDEX idx_created_at (created_at),
	INDEX idx_deleted_at (deleted_at)
);`, tableIdentifier)

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	// Best-effort upgrades for already-existing tables.
	_, _ = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN updated_at TIMESTAMP NULL DEFAULT NULL", tableIdentifier))
	_, _ = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL", tableIdentifier))
	_, _ = db.Exec(fmt.Sprintf("CREATE INDEX idx_deleted_at ON %s (deleted_at)", tableIdentifier))

	return nil
}
