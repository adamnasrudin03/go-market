package models

type Order struct {
	ID        uint64  `json:"id" gorm:"primaryKey"`
	ProductID uint64  `json:"product_id" gorm:"not null"`
	ShopID    uint64  `json:"shop_id" gorm:"not null"`
	Quantity  uint64  `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	DefaultModel
}
