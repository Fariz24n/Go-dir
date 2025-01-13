package models

type Transaction struct {
	No         uint    `json:"no" gorm:"primaryKey"`
	CustomerID uint    `json:"customer_id"`
	Item       string  `json:"item"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}
