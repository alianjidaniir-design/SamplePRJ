package task

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	mysqlDataSource "github.com/alianjidaniir-design/SamplePRJ/models/task/dataSources/mysql"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

type Repository struct {
	idCounter int64
	tasks     []taskDataModel.Task
	lock      sync.RWMutex
	listCache map[string]taskSchema.ListResponse
	cacheLock sync.RWMutex

	db        *sql.DB
	tableName string
	tableSQL  string
	initErr   error
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			idCounter: 100,
			tasks:     []taskDataModel.Task{},
			listCache: map[string]taskSchema.ListResponse{},
		}
		repoIns.initializeStorage()
	})
	return repoIns
}

func init() {
	repositories.TaskRepo = GetRepo()
}

func (repo *Repository) initializeStorage() {
	cfg, err := mysqlDataSource.LoadConfigFromEnv()
	if err != nil {
		repo.initErr = fmt.Errorf("load mysql config failed: %w", err)
		return
	}

	if cfg.DSN == "" {
		log.Println("[task-repository] MYSQL_DSN is empty, using in-memory repository")
		return
	}

	tableSQL, err := mysqlDataSource.TaskTableIdentifier(cfg.TaskTableName)
	if err != nil {
		repo.initErr = err
		return
	}

	db, err := mysqlDataSource.Open(cfg)
	if err != nil {
		repo.initErr = fmt.Errorf("open mysql failed: %w", err)
		return
	}

	if err := mysqlDataSource.EnsureTaskTable(db, cfg.TaskTableName); err != nil {
		_ = db.Close()
		repo.initErr = fmt.Errorf("create task table failed: %w", err)
		return
	}

	repo.db = db
	repo.tableName = cfg.TaskTableName
	repo.tableSQL = tableSQL

	log.Printf("[task-repository] mysql storage enabled table=%s", repo.tableName)
}

func (repo *Repository) nextID() int64 {
	return atomic.AddInt64(&repo.idCounter, 1)
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest]) (res taskSchema.CreateResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return taskSchema.CreateResponse{}, "03", status.StatusInternalServerError, repo.initErr
	}

	if repo.db != nil {
		return repo.createInMySQL(ctx, req)
	}

	return repo.createInMemory(req)
}

func (repo *Repository) createInMemory(req commonSchema.BaseRequest[taskSchema.CreateRequest]) (res taskSchema.CreateResponse, errStr string, code int, err error) {
	task := taskDataModel.Task{
		ID:          repo.nextID(),
		Title:       req.Body.Title,
		Description: req.Body.Description,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	repo.lock.Lock()
	repo.tasks = append(repo.tasks, task)
	repo.lock.Unlock()

	repo.cacheLock.Lock()
	repo.listCache = map[string]taskSchema.ListResponse{}
	repo.cacheLock.Unlock()

	res = taskSchema.CreateResponse{Task: task}
	return res, "", status.StatusOK, nil
}

func (repo *Repository) createInMySQL(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest]) (res taskSchema.CreateResponse, errStr string, code int, err error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", repo.tableSQL)
	insertResult, err := repo.db.ExecContext(ctx, insertQuery, req.Body.Title, req.Body.Description)
	if err != nil {
		return taskSchema.CreateResponse{}, "04", status.StatusInternalServerError, err
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return taskSchema.CreateResponse{}, "05", status.StatusInternalServerError, err
	}

	var task taskDataModel.Task
	var createdAt time.Time
	readQuery := fmt.Sprintf("SELECT id, title, description, created_at FROM %s WHERE id = ?", repo.tableSQL)
	err = repo.db.QueryRowContext(ctx, readQuery, insertedID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&createdAt,
	)
	if err != nil {
		return taskSchema.CreateResponse{}, "06", status.StatusInternalServerError, err
	}

	task.CreatedAt = createdAt.UTC().Format(time.RFC3339)

	repo.clearListCache()

	return taskSchema.CreateResponse{Task: task}, "", status.StatusOK, nil
}

func (repo *Repository) clearListCache() {
	repo.cacheLock.Lock()
	repo.listCache = map[string]taskSchema.ListResponse{}
	repo.cacheLock.Unlock()
}
