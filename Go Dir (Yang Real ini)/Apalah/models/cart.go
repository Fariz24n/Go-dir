package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	No         uint    `json:"no" gorm:"primaryKey"`
	CustomerID uint    `json:"customer_id" gorm:"not null"`
	Item       string  `json:"item" gorm:"type:varchar(255);not null"`
	Quantity   int     `json:"quantity" gorm:"not null"`
	Price      float64 `json:"price" gorm:"type:numeric(6,2);not null"`
	Total      float64 `json:"total" gorm:"type:numeric(6,2);default:0.0"`
}

func MigrateCart(db *gorm.DB) error {
	return db.AutoMigrate(&Cart{})
}
