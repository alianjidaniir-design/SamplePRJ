package memoryDS

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
)

func tehranLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

type TaskDBDS struct {
	idCounter int64
	tasks     []taskRecord
	lock      sync.RWMutex
}

type taskRecord struct {
	task      taskDataModel.Task
	deletedAt *string
}

func NewTaskDBDS(startID int64) *TaskDBDS {
	return &TaskDBDS{
		idCounter: startID,
		tasks:     []taskRecord{},
	}
}

func (ds *TaskDBDS) CreateTask(ctx context.Context, req taskSchema.CreateRequest) (taskDataModel.Task, error) {
	_ = ctx

	now := time.Now().In(tehranLocation())
	task := taskDataModel.Task{
		ID:          atomic.AddInt64(&ds.idCounter, 1),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   now.Format("2006-01-02 15:04:05"),
	}

	ds.lock.Lock()
	ds.tasks = append(ds.tasks, taskRecord{task: task})
	ds.lock.Unlock()

	return task, nil
}

func (ds *TaskDBDS) ListTasks(ctx context.Context, page int, perPage int) ([]taskDataModel.Task, int, error) {
	_ = ctx

	ds.lock.RLock()
	activeTasks := make([]taskDataModel.Task, 0, len(ds.tasks))
	for _, record := range ds.tasks {
		if record.deletedAt != nil {
			continue
		}
		activeTasks = append(activeTasks, record.task)
	}
	ds.lock.RUnlock()

	start := (page - 1) * perPage
	if start > len(activeTasks) {
		start = len(activeTasks)
	}

	end := start + perPage
	if end > len(activeTasks) {
		end = len(activeTasks)
	}

	resultTasks := make([]taskDataModel.Task, end-start)
	copy(resultTasks, activeTasks[start:end])

	return resultTasks, len(activeTasks), nil
}

func (ds *TaskDBDS) UpdateTask(ctx context.Context, req taskSchema.UpdateRequest) (taskDataModel.Task, bool, error) {
	_ = ctx

	now := time.Now().In(tehranLocation()).Format("2006-01-02 15:04:05")

	ds.lock.Lock()
	defer ds.lock.Unlock()

	for i := range ds.tasks {
		if ds.tasks[i].task.ID != req.TaskID {
			continue
		}
		if ds.tasks[i].deletedAt != nil {
			return taskDataModel.Task{}, false, nil
		}

		if req.Title != nil {
			ds.tasks[i].task.Title = *req.Title
		}
		if req.Description != nil {
			ds.tasks[i].task.Description = *req.Description
		}

		ds.tasks[i].task.UpdatedAt = &now
		return ds.tasks[i].task, true, nil
	}

	return taskDataModel.Task{}, false, nil
}

func (ds *TaskDBDS) SoftDeleteTask(ctx context.Context, taskID int64) (taskDataModel.Task, bool, error) {
	_ = ctx

	now := time.Now().In(tehranLocation()).Format("2006-01-02 15:04:05")

	ds.lock.Lock()
	defer ds.lock.Unlock()

	for i := range ds.tasks {
		if ds.tasks[i].task.ID != taskID {
			continue
		}
		if ds.tasks[i].deletedAt != nil {
			return taskDataModel.Task{}, false, nil
		}

		ds.tasks[i].deletedAt = &now
		ds.tasks[i].task.DeletedAt = &now
		ds.tasks[i].task.UpdatedAt = &now
		return ds.tasks[i].task, true, nil
	}

	return taskDataModel.Task{}, false, nil
}

func (ds *TaskDBDS) Reset() {
	ds.lock.Lock()
	ds.tasks = []taskRecord{}
	ds.idCounter = 100
	ds.lock.Unlock()
}
