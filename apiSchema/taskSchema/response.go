package taskSchema

import "github.com/alianjidaniir-design/SamplePRJ/models/datamodel"

type CreateResponse struct {
	Task datamodel.Task `json:"task" msgpack:"task"`
}

type ListResponse struct {
	Tasks   []datamodel.Task `json:"tasks" msgpack:"tasks"`
	Page    int              `json:"page" msgpack:"page"`
	PerPage int              `json:"perPage" msgpack:"perPage"`
	Total   int              `json:"total" msgpack:"total"`
}
