package task

import (
	"context"
	"fmt"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[taskSchema.ListRequest], user datamodel.User) (res taskSchema.ListResponse, errStr string, code int, err error) {
	_ = ctx
	_ = user

	cacheKey := fmt.Sprintf("task:list:page:%d:perPage:%d", req.Body.Page, req.Body.PerPage)

	repo.cacheLock.RLock()
	cachedRes, cacheHit := repo.listCache[cacheKey]
	repo.cacheLock.RUnlock()
	if cacheHit {
		return cloneListResponse(cachedRes), "", status.StatusOK, nil
	}

	repo.lock.RLock()
	clonedTasks := make([]datamodel.Task, len(repo.tasks))
	copy(clonedTasks, repo.tasks)
	repo.lock.RUnlock()

	start := (req.Body.Page - 1) * req.Body.PerPage
	if start > len(clonedTasks) {
		start = len(clonedTasks)
	}

	end := start + req.Body.PerPage
	if end > len(clonedTasks) {
		end = len(clonedTasks)
	}

	res = taskSchema.ListResponse{
		Tasks:   clonedTasks[start:end],
		Page:    req.Body.Page,
		PerPage: req.Body.PerPage,
		Total:   len(clonedTasks),
	}

	repo.cacheLock.Lock()
	repo.listCache[cacheKey] = cloneListResponse(res)
	repo.cacheLock.Unlock()

	return res, "", status.StatusOK, nil
}

func cloneListResponse(source taskSchema.ListResponse) taskSchema.ListResponse {
	cloned := source
	cloned.Tasks = make([]datamodel.Task, len(source.Tasks))
	copy(cloned.Tasks, source.Tasks)
	return cloned
}
