package student

import (
	"Fiber/API/2/apiSchema/commonSchema"
	"Fiber/API/2/apiSchema/studentsSchema"
	"Fiber/API/2/statics/constants/status"
	"context"
	"errors"
)

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[studentsSchema.CreateUserRequest]) (res studentsSchema.UserLoginResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentsSchema.UserLoginResponse{}, "03", status.StatusNotImplemented, repo.initErr
	}
	if repo.db() == nil {
		return studentsSchema.UserLoginResponse{}, "03", status.StatusNotImplemented, errors.New("student datasource not configured")
	}
	createdUser, err := repo.db().CreateStudent(ctx, req.Body)
	if err != nil {
		return studentsSchema.UserLoginResponse{}, "04", status.StatusNotImplemented, err
	}
	return studentsSchema.UserLoginResponse{Students: createdUser}, "", status.StatusOK, nil
}
