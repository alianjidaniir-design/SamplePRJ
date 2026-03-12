package repositories

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
)

type UserRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.CreateRequest]) (res userSchema.CreateResponse, errStr string, code int, err error)
	Info(ctx context.Context, req commonSchema.BaseRequest[userSchema.InfoRequest]) (res userSchema.InfoResponse, errStr string, code int, err error)
}

var UserRepo UserRepository
