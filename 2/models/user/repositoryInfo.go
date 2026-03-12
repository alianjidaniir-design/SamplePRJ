package user

import (
	"context"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
	"github.com/alianjidaniir-design/SamplePRJ/statics/customErr"
)

func (repo *Repository) Info(ctx context.Context, req commonSchema.BaseRequest[userSchema.InfoRequest]) (res userSchema.InfoResponse, errStr string, code int, err error) {
	_ = ctx

	repo.lock.RLock()
	defer repo.lock.RUnlock()

	for _, currentUser := range repo.users {
		if currentUser.ID == req.Body.UserID {
			return userSchema.InfoResponse{User: currentUser}, "", status.StatusOK, nil
		}
	}

	return userSchema.InfoResponse{}, "12", status.StatusBadRequest, customErr.UserNotFound
}
