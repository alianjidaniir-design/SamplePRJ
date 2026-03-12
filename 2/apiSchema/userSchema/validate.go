package userSchema

import (
	"strings"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
	"github.com/alianjidaniir-design/SamplePRJ/statics/customErr"
)

func (req *CreateRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" {
		return "03", status.StatusBadRequest, customErr.InvalidUsername
	}

	if req.Email == "" {
		return "06", status.StatusBadRequest, customErr.InvalidEmail
	}

	_ = validateExtraData
	return "", status.StatusOK, nil
}

func (req *InfoRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error) {
	if req.UserID < 1 {
		return "09", status.StatusBadRequest, customErr.InvalidUserID
	}

	_ = validateExtraData
	return "", status.StatusOK, nil
}
