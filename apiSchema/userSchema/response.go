package userSchema

import "github.com/alianjidaniir-design/SamplePRJ/models/datamodel"

type CreateResponse struct {
	User datamodel.User `json:"user" msgpack:"user"`
}

type InfoResponse struct {
	User datamodel.User `json:"user" msgpack:"user"`
}
