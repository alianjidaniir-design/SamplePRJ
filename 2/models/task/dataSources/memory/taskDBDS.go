package memory

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
)

type TaskDBDS struct {
	idCounter int64
	tasks     []taskDataModel.Task
	lock      sync.RWMutex
}

func NewTaskDBDS(startID int64) *TaskDBDS {
	return &TaskDBDS{
		idCounter: startID,
		tasks:     []taskDataModel.Task{},
	}
}

func (ds *TaskDBDS) CreateTask(ctx context.Context, req taskSchema.CreateRequest) (taskDataModel.Task, error) {
	_ = ctx

	task := taskDataModel.Task{
		ID:          atomic.AddInt64(&ds.idCounter, 1),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	ds.lock.Lock()
	ds.tasks = append(ds.tasks, task)
	ds.lock.Unlock()

	return task, nil
}

func (ds *TaskDBDS) ListTasks(ctx context.Context, page int, perPage int) ([]taskDataModel.Task, int, error) {
	_ = ctx

	ds.lock.RLock()
	clonedTasks := make([]taskDataModel.Task, len(ds.tasks))
	copy(clonedTasks, ds.tasks)
	ds.lock.RUnlock()

	start := (page - 1) * perPage
	if start > len(clonedTasks) {
		start = len(clonedTasks)
	}

	end := start + perPage
	if end > len(clonedTasks) {
		end = len(clonedTasks)
	}

	resultTasks := make([]taskDataModel.Task, end-start)
	copy(resultTasks, clonedTasks[start:end])

	return resultTasks, len(clonedTasks), nil
}

func (ds *TaskDBDS) Reset() {
	ds.lock.Lock()
	ds.tasks = []taskDataModel.Task{}
	ds.idCounter = 100
	ds.lock.Unlock()
}
