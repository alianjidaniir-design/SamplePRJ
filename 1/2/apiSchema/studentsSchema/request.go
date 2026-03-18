package studentsSchema

type CreateUserRequest struct {
	StudentCode string `json:"student_code" msgpack:"student_code" validate:"required,max=128"`
	FirstName   string `json:"first_name" msgpack:"first_name" validate:"required"`
	LastName    string `json:"last_name" msgpack:"last_name" validate:"required"`
}
