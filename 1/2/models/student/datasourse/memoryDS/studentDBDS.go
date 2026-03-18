package memoryDS

import (
	"Fiber/API/2/apiSchema/studentsSchema"
	studentDataModel "Fiber/API/2/models/student/dataModel"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type StudentDBDS struct {
	idCounter int64
	students  []studentDataModel.Students
	lock      sync.RWMutex
}

func NewStudentDBDS(startID int64) *StudentDBDS {
	return &StudentDBDS{
		idCounter: startID,
		students:  []studentDataModel.Students{},
	}
}

func (db *StudentDBDS) CreateStudent(ctx context.Context, req studentsSchema.CreateUserRequest) (studentDataModel.Students, error) {
	_ = ctx

	student := studentDataModel.Students{
		ID:          atomic.AddInt64(&db.idCounter, 1),
		StudentCode: req.StudentCode,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	db.lock.Lock()
	db.students = append(db.students, student)
	db.lock.Unlock()
	return student, nil

}
