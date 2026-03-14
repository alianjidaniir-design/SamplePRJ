package dataModel

type Task struct {
	ID          int64   `gorm:"column:id;primaryKey" json:"id" msgpack:"id"`
	Title       string  `gorm:"column:title" json:"title" msgpack:"title"`
	Description string  `gorm:"column:description" json:"description" msgpack:"description"`
	CreatedAt   string  `gorm:"column:created_at" json:"createdAt" msgpack:"createdAt"`
	UpdatedAt   *string `gorm:"column:updated_at" json:"updatedAt,omitempty" msgpack:"updatedAt,omitempty"`
	DeletedAt   *string `gorm:"column:deleted_at" json:"deletedAt,omitempty" msgpack:"deletedAt,omitempty"`
}
