package models

type Warehouse struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	DefaultModel
}
