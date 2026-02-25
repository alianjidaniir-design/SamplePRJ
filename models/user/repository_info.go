package user

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
	"github.com/alianjidaniir-design/SamplePRJ/statics/customErr"
)

func (repo *Repository) Info(ctx context.Context, req commonSchema.BaseRequest[userSchema.InfoRequest], user datamodel.User) (res userSchema.InfoResponse, errStr string, code int, err error) {
	_ = ctx
	_ = user

	repo.lock.RLock()
	defer repo.lock.RUnlock()

	for _, currentUser := range repo.users {
		if currentUser.ID == req.Body.UserID {
			res = userSchema.InfoResponse{User: currentUser}
			return res, "", status.StatusOK, nil
		}
	}

	return userSchema.InfoResponse{}, "12", status.StatusBadRequest, customErr.UserNotFound
}
