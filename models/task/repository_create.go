package task

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

type Repository struct {
	idCounter int64
	tasks     []taskDataModel.Task
	lock      sync.RWMutex
	listCache map[string]taskSchema.ListResponse
	cacheLock sync.RWMutex
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
	})
	return repoIns
}

func init() {
	repositories.TaskRepo = GetRepo()
}

func (repo *Repository) nextID() int64 {
	return atomic.AddInt64(&repo.idCounter, 1)
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest]) (res taskSchema.CreateResponse, errStr string, code int, err error) {
	_ = ctx

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
