package repositories

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
)

type TaskRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest], user datamodel.User) (res taskSchema.CreateResponse, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[taskSchema.ListRequest], user datamodel.User) (res taskSchema.ListResponse, errStr string, code int, err error)
}

var TaskRepo TaskRepository
