package dataSources

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
)

type TaskDBDS interface {
	CreateTask(ctx context.Context, req taskSchema.CreateRequest) (taskDataModel.Task, error)
	ListTasks(ctx context.Context, page int, perPage int) ([]taskDataModel.Task, int, error)
}

type TaskCacheDS interface {
	GetList(cacheKey string) (taskSchema.ListResponse, bool)
	SetList(cacheKey string, res taskSchema.ListResponse)
	InvalidateList()
}
