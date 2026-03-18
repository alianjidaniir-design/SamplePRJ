package studentsSchema

import (
	"Fiber/API/2/apiSchema/commonSchema"
	"Fiber/API/2/statics/constants/status"
	"Fiber/API/2/statics/customErr"
	"fmt"
	"strings"
)

func (req *CreateUserRequest) Validate(validataExtraData commonSchema.ValidateExtraData) (string, int, error) {
	req.StudentCode = strings.TrimSpace(req.StudentCode)
	fmt.Println(req.StudentCode)
	if req.StudentCode == "" {
		return "03", status.StatusBadRequest, customErr.InvalidStudentCode
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	if req.FirstName == "" {
		return "06", status.StatusBadRequest, customErr.InvalidFirstName
	}
	req.LastName = strings.TrimSpace(req.LastName)
	if req.LastName == "" {
		return "09", status.StatusBadRequest, customErr.InvalidLastName
	}

	_ = validataExtraData
	return "", status.StatusOK, nil

}
