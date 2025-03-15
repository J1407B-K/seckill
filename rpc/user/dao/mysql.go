package dao

import (
	"gorm.io/gorm"
	"seckill/idl/kitex_gen/user"
	"seckill/rpc/user/model"
)

func SaveUser(db *gorm.DB, req *user.RegisterReq) error {
	var user model.User

	user.Username = req.Username
	user.Password = req.Password
	user.Email = req.Email

	db.Create(&user)
	return nil
}
