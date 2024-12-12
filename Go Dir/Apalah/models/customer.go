package models

type Customer struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}
