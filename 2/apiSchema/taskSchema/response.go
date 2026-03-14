package taskSchema

import (
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/dataModel"
)

type CreateResponse struct {
	Task taskDataModel.Task `json:"task" msgpack:"task"`
}

type ListResponse struct {
	Tasks   []taskDataModel.Task `json:"tasks" msgpack:"tasks"`
	Page    int                  `json:"page" msgpack:"page"`
	PerPage int                  `json:"perPage" msgpack:"perPage"`
	Total   int                  `json:"total" msgpack:"total"`
}

type UpdateResponse struct {
	Task taskDataModel.Task `json:"task" msgpack:"task"`
}

type DeleteResponse struct {
	Task taskDataModel.Task `json:"task" msgpack:"task"`
}
