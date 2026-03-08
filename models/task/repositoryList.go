package task

import (
	"context"
	"fmt"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[taskSchema.ListRequest]) (res taskSchema.ListResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return taskSchema.ListResponse{}, "03", status.StatusInternalServerError, repo.initErr
	}

	cacheKey := fmt.Sprintf("task:list:page:%d:perPage:%d", req.Body.Page, req.Body.PerPage)
	cachedRes, cacheHit := repo.cache().GetList(cacheKey)
	if cacheHit {
		return cloneListResponse(cachedRes), "", status.StatusOK, nil
	}

	tasks, total, err := repo.db().ListTasks(ctx, req.Body.Page, req.Body.PerPage)
	if err != nil {
		return taskSchema.ListResponse{}, "04", status.StatusInternalServerError, err
	}

	res = taskSchema.ListResponse{
		Tasks:   tasks,
		Page:    req.Body.Page,
		PerPage: req.Body.PerPage,
		Total:   total,
	}

	repo.cache().SetList(cacheKey, cloneListResponse(res))
	return res, "", status.StatusOK, nil
}

func cloneListResponse(source taskSchema.ListResponse) taskSchema.ListResponse {
	cloned := source
	cloned.Tasks = make([]taskDataModel.Task, len(source.Tasks))
	copy(cloned.Tasks, source.Tasks)
	return cloned
}
