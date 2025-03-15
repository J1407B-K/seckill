package flag

import (
	"gorm.io/gorm"
	"log"
	"seckill/rpc/user/model"
)

func MysqlAutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&model.User{},
		)
	if err != nil {
		log.Println("建表失败")
	}
	log.Println("建表成功")
}
