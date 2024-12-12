package models

type Aisle struct {
	No       uint    `json:"no" gorm:"primaryKey"`
	Supplier string  `json:"supplier" gorm:"not null"`
	Item     string  `json:"item" gorm:"not null"`
	Quantity int     `json:"quantity" gorm:"not null"`
	Ket      string  `json:"ket" gorm:"not null"`
	Price    float64 `json:"price" gorm:"not null"`
}
