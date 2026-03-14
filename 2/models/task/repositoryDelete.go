package task

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
	"github.com/alianjidaniir-design/SamplePRJ/statics/customErr"
)

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[taskSchema.DeleteRequest]) (res taskSchema.DeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return taskSchema.DeleteResponse{}, "03", status.StatusInternalServerError, repo.initErr
	}

	deletedTask, found, err := repo.db().SoftDeleteTask(ctx, req.Body.TaskID)
	if err != nil {
		return taskSchema.DeleteResponse{}, "04", status.StatusInternalServerError, err
	}
	if !found {
		return taskSchema.DeleteResponse{}, "12", status.StatusBadRequest, customErr.TaskNotFound
	}

	repo.cache().InvalidateList()
	return taskSchema.DeleteResponse{Task: deletedTask}, "", status.StatusOK, nil
}
