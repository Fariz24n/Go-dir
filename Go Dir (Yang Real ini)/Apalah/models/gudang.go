package models

import (
	"gorm.io/gorm"
)

type Gudang struct {
	No       uint    `json:"no" gorm:"primaryKey"`
	Supplier string  `json:"supplier" gorm:"type:varchar(255);not null"`
	Item     string  `json:"item" gorm:"type:varchar(255);not null"`
	Quantity int     `json:"quantity" gorm:"not null"`
	Ket      string  `json:"ket" gorm:"type:varchar(255)"`
	Price    float64 `json:"price" gorm:"not null"`
}

func MigrateGudang(db *gorm.DB) error {
	return db.AutoMigrate(&Gudang{})
}
