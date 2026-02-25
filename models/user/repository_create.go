package user

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.CreateRequest], user datamodel.User) (res userSchema.CreateResponse, errStr string, code int, err error) {
	_ = ctx
	_ = user

	newUser := datamodel.User{
		ID:       repo.nextID(),
		Username: req.Body.Username,
		Email:    req.Body.Email,
	}

	repo.lock.Lock()
	repo.users = append(repo.users, newUser)
	repo.lock.Unlock()

	res = userSchema.CreateResponse{User: newUser}
	return res, "", status.StatusOK, nil
}
