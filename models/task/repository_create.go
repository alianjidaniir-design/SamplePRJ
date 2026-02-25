package task

import (
	"context"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest], user datamodel.User) (res taskSchema.CreateResponse, errStr string, code int, err error) {
	_ = ctx
	_ = user

	task := datamodel.Task{
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
