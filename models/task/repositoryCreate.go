package task

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest]) (res taskSchema.CreateResponse, errStr string, code int, err error) {
	createdTask, err := repo.db().CreateTask(ctx, req.Body)
	if err != nil {
		return taskSchema.CreateResponse{}, "04", status.StatusInternalServerError, err
	}

	repo.cache().InvalidateList()
	return taskSchema.CreateResponse{Task: createdTask}, "", status.StatusOK, nil
}
