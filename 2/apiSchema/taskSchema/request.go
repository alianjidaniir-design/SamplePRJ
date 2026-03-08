package taskSchema

type CreateRequest struct {
	Title       string `json:"title" msgpack:"title" validate:"required,max=128"`
	Description string `json:"description" msgpack:"description" validate:"max=512"`
}

type ListRequest struct {
	Page    int `json:"page" msgpack:"page" validate:"required"`
	PerPage int `json:"perPage" msgpack:"perPage" validate:"required"`
}
