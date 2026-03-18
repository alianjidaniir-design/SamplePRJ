package studentsSchema

import "Fiber/API/2/models/student/dataModel"

type UserLoginResponse struct {
	Students dataModel.Students `json:"user" msgpack:"user" `
}
