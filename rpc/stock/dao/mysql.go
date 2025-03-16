package dao

import (
	"gorm.io/gorm"
	"seckill/rpc/stock/model"
)

func MysqlSearchStock(db *gorm.DB, id int) *model.ProductStock {
	var stock *model.ProductStock

	err := db.Where("product_id = ?", id).First(&stock).Error
	if err != nil {
		return nil
	}
	return stock
}
