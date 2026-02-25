package datamodel

type Task struct {
	ID          int64  `json:"id" msgpack:"id"`
	Title       string `json:"title" msgpack:"title"`
	Description string `json:"description" msgpack:"description"`
	CreatedAt   string `json:"createdAt" msgpack:"createdAt"`
}

type User struct {
	ID       int64  `json:"id" msgpack:"id"`
	Username string `json:"username" msgpack:"username"`
	Email    string `json:"email" msgpack:"email"`
}
