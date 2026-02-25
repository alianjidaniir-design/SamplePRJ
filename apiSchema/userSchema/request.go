package userSchema

type CreateRequest struct {
	Username string `json:"username" msgpack:"username" validate:"required,max=64"`
	Email    string `json:"email" msgpack:"email" validate:"required,max=128"`
}

type InfoRequest struct {
	UserID int64 `json:"userID" msgpack:"userID" validate:"required"`
}
