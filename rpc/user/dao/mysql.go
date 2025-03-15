package dao

import (
	"gorm.io/gorm"
	"seckill/idl/kitex_gen/user"
	"seckill/rpc/user/model"
)

func SaveUser(db *gorm.DB, req *user.RegisterReq) error {
	// 创建用户
	u := model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := db.Create(&u).Error; err != nil {
		return err
	}
	return nil
}
