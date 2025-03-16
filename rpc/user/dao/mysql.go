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

func SelectUser(db *gorm.DB, k string) (*model.User, error) {
	var u model.User

	err := db.Where("username = ?", k).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
