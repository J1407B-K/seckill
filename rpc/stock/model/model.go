package model

type ProductStock struct {
	ProductId int   `json:"productId" gorm:"primaryKey"`
	Stock     int32 `json:"stock" gorm:"not null"`
}
