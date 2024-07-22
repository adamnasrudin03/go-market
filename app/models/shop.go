package models

type Shop struct {
	ID     uint64 `json:"id" gorm:"primary_key;auto_increment"`
	UserID uint64 `json:"user_id" gorm:"not null"`
	Name   string `json:"name" gorm:"not null;unique"`
	Status string `json:"status" gorm:"not null"`
	DefaultModel
}
