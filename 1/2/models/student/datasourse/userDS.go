package datasourse

import (
	"Fiber/API/2/apiSchema/studentsSchema"
	"context"

	studentDataModel "Fiber/API/2/models/student/dataModel"
)

type StudentDBDS interface {
	CreateStudent(ctx context.Context, req studentsSchema.CreateUserRequest) (studentDataModel.Students, error)
}
