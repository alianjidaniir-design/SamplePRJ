package repositories

import (
	"Fiber/API/2/apiSchema/commonSchema"
	"Fiber/API/2/apiSchema/studentsSchema"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[studentsSchema.CreateUserRequest]) (res studentsSchema.UserLoginResponse, errStr string, code int, err error)
}

var UserRepo UserRepository
