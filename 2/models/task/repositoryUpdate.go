package task

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
	"github.com/alianjidaniir-design/SamplePRJ/statics/customErr"
)

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[taskSchema.UpdateRequest]) (res taskSchema.UpdateResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return taskSchema.UpdateResponse{}, "03", status.StatusInternalServerError, repo.initErr
	}

	updatedTask, found, err := repo.db().UpdateTask(ctx, req.Body)
	if err != nil {
		return taskSchema.UpdateResponse{}, "04", status.StatusInternalServerError, err
	}
	if !found {
		return taskSchema.UpdateResponse{}, "12", status.StatusBadRequest, customErr.TaskNotFound
	}

	repo.cache().InvalidateList()
	return taskSchema.UpdateResponse{Task: updatedTask}, "", status.StatusOK, nil
}
