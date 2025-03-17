package main

import (
	"context"
	"gorm.io/gorm"
	"log"
	"seckill/idl/kitex_gen/common"
	user "seckill/idl/kitex_gen/user"
	"seckill/rpc/user/dao"
	"seckill/rpc/user/hash"
	"strconv"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	DB *gorm.DB
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	req.Password, err = hash.HashedLock(req.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = dao.SaveUser(s.DB, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user.RegisterResp{
		Resp: &common.Resp{
			Code: 0,
			Msg:  "ok",
			Data: req.Username,
		},
	}, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	//获得userinfo
	userinfo, err := dao.SelectUser(s.DB, req.Username)
	if err != nil {
		return nil, err
	}

	//比较是否一致
	err = hash.CompareHashAndPassword(userinfo.Password, req.Password)
	if err != nil {
		return nil, err
	}

	resp = &user.LoginResp{
		Resp: &common.Resp{
			Code: 0,
			Msg:  "ok",
			Data: strconv.Itoa(userinfo.Userid),
		},
	}
	return resp, nil
}
