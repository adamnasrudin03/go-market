package seeders

import (
	"github.com/adamnasrudin03/go-market/app/models"
	"gorm.io/gorm"
)

func InitProducts(db *gorm.DB) {
	tx := db.Begin()
	var warehouse = []models.Warehouse{}
	tx.Select("id").Limit(2).Find(&warehouse)
	if len(warehouse) == 0 {
		warehouse = []models.Warehouse{
			{
				Name: "Warehouse A",
			},
			{
				Name: "Warehouse B",
			},
		}
		tx.Create(&warehouse)
	}

	var products = []models.Product{}
	tx.Select("id").Limit(1).Find(&products)
	if len(products) == 0 {
		products = []models.Product{
			{
				Name:        "Product 1",
				WarehouseID: warehouse[0].ID,
				Stock:       10,
				Price:       10000,
			},
			{
				Name:        "Product 2",
				WarehouseID: warehouse[0].ID,
				Stock:       20,
				Price:       10500,
			},
			{
				Name:        "Product 3",
				WarehouseID: warehouse[0].ID,
				Stock:       30,
				Price:       2500,
			},
			{
				Name:        "Product 4",
				WarehouseID: warehouse[1].ID,
				Stock:       30,
				Price:       2500,
			},
			{
				Name:        "Product 5",
				WarehouseID: warehouse[1].ID,
				Stock:       30,
				Price:       10500,
			},
		}
		tx.Create(&products)
	}

	tx.Commit()
}
