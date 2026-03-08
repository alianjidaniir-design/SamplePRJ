package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
)

type TaskDBDS struct {
	tableName string
	tableSQL  string
	db        DBExecutor
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

	var task taskDataModel.Task
	var createdAt time.Time
	readQuery := fmt.Sprintf("SELECT id, title, description, created_at FROM %s WHERE id = ?", ds.tableSQL)
	err = ds.db.QueryRowContext(ctx, readQuery, insertedID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&createdAt,
	)
	if err != nil {
		return taskDataModel.Task{}, err
	}

	task.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return task, nil
}

func (ds *TaskDBDS) ListTasks(ctx context.Context, page int, perPage int) ([]taskDataModel.Task, int, error) {
	offset := (page - 1) * perPage

	rowsQuery := fmt.Sprintf(
		"SELECT id, title, description, created_at FROM %s ORDER BY id ASC LIMIT ? OFFSET ?",
		ds.tableSQL,
	)
	rows, err := ds.db.QueryContext(ctx, rowsQuery, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]taskDataModel.Task, 0, perPage)
	for rows.Next() {
		var task taskDataModel.Task
		var createdAt time.Time
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &createdAt); err != nil {
			return nil, 0, err
		}

		task.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableSQL)
	total := 0
	if err := ds.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (ds *TaskDBDS) TableName() string {
	return ds.tableName
}
