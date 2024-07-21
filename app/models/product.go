package models

type Product struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"not null"`
	Stock       uint64  `json:"stock" gorm:"not null"`
	Price       float64 `json:"price"`
	WarehouseID uint64  `json:"warehouse_id" gorm:"not null"`
	DefaultModel
}
