package model

import "time"

type Order struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`     // 自增主键
	OrderID   string    `gorm:"uniqueIndex;size:64;not null"` // 订单号，唯一标识
	UserID    string    `gorm:"size:64;not null"`             // 用户
	ProductID string    `gorm:"size:64;not null"`             // 商品ID
	Count     int       `gorm:"not null"`                     // 购买数量
	Status    int       `gorm:"not null"`                     // 订单状态：0-待支付，1-已支付，2-已取消，3-回滚
	CreatedAt time.Time `gorm:"autoCreateTime"`               // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`               // 更新时间
}
