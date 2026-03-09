package customErr

import (
	"errors"

	"github.com/alianjidaniir-design/SamplePRJ/statics/constants"
)

var (
	InvalidTitle   = errors.New(constants.InvalidTitle)
	InvalidPage    = errors.New(constants.InvalidPage)
	InvalidPerPage = errors.New(constants.InvalidPerPage)
)
