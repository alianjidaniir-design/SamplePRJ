package mysqlDS

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
	"github.com/go-sql-driver/mysql"
)

type TaskDBDS struct {
	tableName string
	tableSQL  string
	db        DBExecutor
}

func tehranLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func isUnknownColumnErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1054
}

func NewTaskDBDSFromEnv() (*TaskDBDS, bool, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, false, fmt.Errorf("load mysql config failed: %w", err)
	}

	if cfg.DSN == "" {
		return nil, false, nil
	}

	tableSQL, err := TaskTableIdentifier(cfg.TaskTableName)
	if err != nil {
		return nil, false, err
	}

	db, err := Open(cfg)
	if err != nil {
		return nil, false, fmt.Errorf("open mysql failed: %w", err)
	}

	if err := EnsureTaskTable(db, cfg.TaskTableName); err != nil {
		_ = db.Close()
		return nil, false, fmt.Errorf("create task table failed: %w", err)
	}

	return &TaskDBDS{
		tableName: cfg.TaskTableName,
		tableSQL:  tableSQL,
		db:        db,
	}, true, nil
}

func (ds *TaskDBDS) CreateTask(ctx context.Context, req taskSchema.CreateRequest) (taskDataModel.Task, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", ds.tableSQL)
	insertResult, err := ds.db.ExecContext(ctx, insertQuery, req.Title, req.Description)
	if err != nil {
		return taskDataModel.Task{}, err
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return taskDataModel.Task{}, err
	}

	return ds.readTaskByID(ctx, insertedID)
}

func (ds *TaskDBDS) ListTasks(ctx context.Context, page int, perPage int) ([]taskDataModel.Task, int, error) {
	offset := (page - 1) * perPage

	createdColumn := "created_at"
	rowsQuery := fmt.Sprintf(
		"SELECT id, title, description, %s, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL ORDER BY id ASC LIMIT ? OFFSET ?",
		createdColumn,
		ds.tableSQL,
	)
	rows, err := ds.db.QueryContext(ctx, rowsQuery, perPage, offset)
	if err != nil && isUnknownColumnErr(err) {
		createdColumn = "createdAt"
		rowsQuery = fmt.Sprintf(
			"SELECT id, title, description, %s FROM %s ORDER BY id ASC LIMIT ? OFFSET ?",
			createdColumn,
			ds.tableSQL,
		)
		rows, err = ds.db.QueryContext(ctx, rowsQuery, perPage, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]taskDataModel.Task, 0, perPage)
	for rows.Next() {
		var task taskDataModel.Task
		var createdAt time.Time
		var updatedAt sql.NullTime
		var deletedAt sql.NullTime

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			var fallbackCreatedAt time.Time
			if err := rows.Scan(&task.ID, &task.Title, &task.Description, &fallbackCreatedAt); err != nil {
				return nil, 0, err
			}
			createdAt = fallbackCreatedAt
		}

		task.CreatedAt = createdAt.In(tehranLocation()).Format("2006-01-02 15:04:05")
		if updatedAt.Valid {
			value := updatedAt.Time.In(tehranLocation()).Format("2006-01-02 15:04:05")
			task.UpdatedAt = &value
		}
		if deletedAt.Valid {
			value := deletedAt.Time.In(tehranLocation()).Format("2006-01-02 15:04:05")
			task.DeletedAt = &value
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL", ds.tableSQL)
	total := 0
	if err := ds.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil && isUnknownColumnErr(err) {
		countQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableSQL)
		err = ds.db.QueryRowContext(ctx, countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (ds *TaskDBDS) UpdateTask(ctx context.Context, req taskSchema.UpdateRequest) (taskDataModel.Task, bool, error) {
	setParts := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if req.Title != nil {
		setParts = append(setParts, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, req.TaskID)

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = ? AND deleted_at IS NULL", ds.tableSQL, joinCSV(setParts))
	result, err := ds.db.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return taskDataModel.Task{}, false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return taskDataModel.Task{}, false, err
	}
	if affected == 0 {
		return taskDataModel.Task{}, false, nil
	}

	task, err := ds.readTaskByID(ctx, req.TaskID)
	if err != nil {
		return taskDataModel.Task{}, false, err
	}

	return task, true, nil
}

func (ds *TaskDBDS) SoftDeleteTask(ctx context.Context, taskID int64) (taskDataModel.Task, bool, error) {
	deleteQuery := fmt.Sprintf("UPDATE %s SET deleted_at = NOW(), updated_at = NOW() WHERE id = ? AND deleted_at IS NULL", ds.tableSQL)
	result, err := ds.db.ExecContext(ctx, deleteQuery, taskID)
	if err != nil {
		return taskDataModel.Task{}, false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return taskDataModel.Task{}, false, err
	}
	if affected == 0 {
		return taskDataModel.Task{}, false, nil
	}

	task, err := ds.readTaskByID(ctx, taskID)
	if err != nil {
		return taskDataModel.Task{}, false, err
	}

	return task, true, nil
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

func (ds *TaskDBDS) readTaskByID(ctx context.Context, taskID int64) (taskDataModel.Task, error) {
	var task taskDataModel.Task
	var createdAt time.Time
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	createdColumn := "created_at"
	readQuery := fmt.Sprintf(
		"SELECT id, title, description, %s, updated_at, deleted_at FROM %s WHERE id = ?",
		createdColumn,
		ds.tableSQL,
	)

	err := ds.db.QueryRowContext(ctx, readQuery, taskID).Scan(&task.ID, &task.Title, &task.Description, &createdAt, &updatedAt, &deletedAt)
	if err != nil && isUnknownColumnErr(err) {
		createdColumn = "createdAt"
		readQuery = fmt.Sprintf("SELECT id, title, description, %s FROM %s WHERE id = ?", createdColumn, ds.tableSQL)
		err = ds.db.QueryRowContext(ctx, readQuery, taskID).Scan(&task.ID, &task.Title, &task.Description, &createdAt)
	}
	if err != nil {
		return taskDataModel.Task{}, err
	}

	task.CreatedAt = createdAt.In(tehranLocation()).Format("2006-01-02 15:04:05")
	if updatedAt.Valid {
		value := updatedAt.Time.In(tehranLocation()).Format("2006-01-02 15:04:05")
		task.UpdatedAt = &value
	}
	if deletedAt.Valid {
		value := deletedAt.Time.In(tehranLocation()).Format("2006-01-02 15:04:05")
		task.DeletedAt = &value
	}

	return task, nil
}

func (ds *TaskDBDS) TableName() string {
	return ds.tableName
}
