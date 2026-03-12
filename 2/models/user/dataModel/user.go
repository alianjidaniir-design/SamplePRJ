package dataModel

type User struct {
	ID       int64  `json:"id" msgpack:"id"`
	Username string `json:"username" msgpack:"username"`
	Email    string `json:"email" msgpack:"email"`
}
