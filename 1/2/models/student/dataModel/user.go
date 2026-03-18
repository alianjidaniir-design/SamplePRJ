package dataModel

type Students struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id" msgpack:"id"`
	StudentCode string `gorm:"column:student_code" json:"student_code" msgpack:"student_code"`
	FirstName   string `gorm:"column:first_name" json:"first_name" msgpack:"first_name"`
	LastName    string `gorm:"column:last_name" json:"last_name" msgpack:"last_name"`
	CreatedAt   string `gorm:"column:created_at" json:"created_at" msgpack:"created_at"`
	UpdatedAt   string `gorm:"column:updated_at" json:"updated_at" msgpack:"updated_at"`
	DeletedAt   string `gorm:"column:deleted_at" json:"deleted_at" msgpack:"deleted_at"`
}
