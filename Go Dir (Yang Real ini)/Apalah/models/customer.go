package models

type Customer struct {
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Address  string `json:"addresss" gorm:"not null"`
	No       uint   `json:"no" gorm:"not null"`
	ID       uint   `json:"id" gorm:"primaryKey"`
	IsAdmin  bool   `gorm:"default:false"`
}
