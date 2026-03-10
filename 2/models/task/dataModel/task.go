package dataModel

type Task struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id" msgpack:"id"`
	Title       string `gorm:"column:title" json:"title" msgpack:"title"`
	Description string `gorm:"column:description" json:"description" msgpack:"description"`
	CreatedAt   string `gorm:"column:createdAt" json:"createdAt" msgpack:"createdAt"`
}
