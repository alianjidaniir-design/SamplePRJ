package mySqlDS

import (
	"Fiber/API/2/apiSchema/studentsSchema"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	studentDataModel "Fiber/API/2/models/student/dataModel"

	"github.com/go-sql-driver/mysql"
)

type UserDBDS struct {
	tablename string
	tableSQL  string
	db        DBExecutor
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func isUnknownColmunErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1054
}

func NewTaskDBDSFromEnv() (*UserDBDS, bool, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, false, err
	}
	if cfg.DSN == "" {
		return nil, false, nil
	}

	tableSQL, err := studentTableIdentifier(cfg.StudentTableName)
	if err != nil {
		return nil, false, err
	}

	db, err := Open(cfg)
	if err != nil {
		return nil, false, err
	}

	if err := EnsureTaskTable(db, cfg.StudentTableName); err != nil {
		return nil, false, err
	}
	return &UserDBDS{
		tablename: cfg.StudentTableName,
		tableSQL:  tableSQL,
		db:        db,
	}, true, nil

}

func (ds *UserDBDS) CreateStudent(ctx context.Context, req studentsSchema.CreateUserRequest) (studentDataModel.Students, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (student_code , first_name , last_name) VALUES (?, ? , ?)", ds.tableSQL)
	insertResult, err := ds.db.ExecContext(ctx, insertQuery, req.StudentCode, req.FirstName, req.LastName)
	if err != nil {
		return studentDataModel.Students{}, err
	}
	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return studentDataModel.Students{}, err
	}
	return ds.readTaskByID(ctx, insertedID)
}

func joinCSV(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	joined := parts[0]
	for i := 1; i < len(parts); i++ {
		joined += ", " + parts[i]
	}
	return joined
}

func (ds *UserDBDS) readTaskByID(ctx context.Context, userID int64) (studentDataModel.Students, error) {
	var students studentDataModel.Students
	var createdAt time.Time
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime
	readQuery := fmt.Sprintf("SELECT id , student_code , first_name , last_name , updated_at , created_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
	if err := ds.db.QueryRowContext(ctx, readQuery, userID).Scan(&students.ID, &students.StudentCode, students.FirstName, students.LastName, &createdAt, &updatedAt, &deletedAt); err != nil {
		return studentDataModel.Students{}, err
	}

	students.CreatedAt = createdAt.In(myLocation()).Format("2006-01-02 15:04:05")
	if updatedAt.Valid {
		value := updatedAt.Time.In(myLocation()).Format("2006-01-02 15:04:05")
		students.UpdatedAt = value
	}
	if deletedAt.Valid {
		value := deletedAt.Time.In(myLocation()).Format("2006-01-02 15:04:05")
		students.DeletedAt = value
	}

	return students, nil
}

func (ds *UserDBDS) TableName() string {
	return ds.tablename
}
