package repositories

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
)

type TaskRepository interface {
	// Create Method
	Create(ctx context.Context, req commonSchema.BaseRequest[taskSchema.CreateRequest]) (res taskSchema.CreateResponse, errStr string, code int, err error)
	// List method
	List(ctx context.Context, req commonSchema.BaseRequest[taskSchema.ListRequest]) (res taskSchema.ListResponse, errStr string, code int, err error)
	// Update method
	Update(ctx context.Context, req commonSchema.BaseRequest[taskSchema.UpdateRequest]) (res taskSchema.UpdateResponse, errStr string, code int, err error)
	// SoftDelete method
	Delete(ctx context.Context, req commonSchema.BaseRequest[taskSchema.DeleteRequest]) (res taskSchema.DeleteResponse, errStr string, code int, err error)
}

var TaskRepo TaskRepository
