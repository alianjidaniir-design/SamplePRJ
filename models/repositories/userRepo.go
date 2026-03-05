package repositories

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	userDataModel "github.com/alianjidaniir-design/SamplePRJ/models/user/datamodel"
)

type UserRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.CreateRequest], user userDataModel.User) (res userSchema.CreateResponse, errStr string, code int, err error)
	Info(ctx context.Context, req commonSchema.BaseRequest[userSchema.InfoRequest], user userDataModel.User) (res userSchema.InfoResponse, errStr string, code int, err error)
}

var UserRepo UserRepository
