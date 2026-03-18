package customErr

import (
	"Fiber/API/2/statics/constants"
	"errors"
)

var (
	InvalidStudentCode = errors.New(constants.InvalidStudentCode)
	InvalidFirstName   = errors.New(constants.InvalidFirstName)
	InvalidLastName    = errors.New(constants.InvalidLastName)
	InvalidEmail       = errors.New(constants.InvalidEmail)
)
