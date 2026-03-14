package customErr

import (
	"errors"

	"github.com/alianjidaniir-design/SamplePRJ/statics/constants"
)

var (
	InvalidTitle    = errors.New(constants.InvalidTitle)
	InvalidPage     = errors.New(constants.InvalidPage)
	InvalidPerPage  = errors.New(constants.InvalidPerPage)
	InvalidTaskID   = errors.New(constants.InvalidTaskID)
	InvalidUpdate   = errors.New(constants.InvalidUpdate)
	TaskNotFound    = errors.New(constants.TaskNotFound)
	InvalidUsername = errors.New(constants.InvalidUsername)
	InvalidEmail    = errors.New(constants.InvalidEmail)
	InvalidUserID   = errors.New(constants.InvalidUserID)
	UserNotFound    = errors.New(constants.UserNotFound)
)
