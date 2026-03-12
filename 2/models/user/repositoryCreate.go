package user

import (
	"context"
	"sync/atomic"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	userDataModel "github.com/alianjidaniir-design/SamplePRJ/models/user/dataModel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) nextID() int64 {
	return atomic.AddInt64(&repo.idCounter, 1)
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.CreateRequest]) (res userSchema.CreateResponse, errStr string, code int, err error) {
	_ = ctx

	newUser := userDataModel.User{
		ID:       repo.nextID(),
		Username: req.Body.Username,
		Email:    req.Body.Email,
	}

	repo.lock.Lock()
	repo.users = append(repo.users, newUser)
	repo.lock.Unlock()

	return userSchema.CreateResponse{User: newUser}, "", status.StatusOK, nil
}
