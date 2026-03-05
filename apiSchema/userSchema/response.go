package userSchema

import userDataModel "github.com/alianjidaniir-design/SamplePRJ/models/user/datamodel"

type CreateResponse struct {
	User userDataModel.User `json:"user" msgpack:"user"`
}

type InfoResponse struct {
	User userDataModel.User `json:"user" msgpack:"user"`
}
