package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"seckill/global"
	userrpc "seckill/idl/kitex_gen/user"
	"seckill/model"
	my_utils "seckill/utils"
)

func Register(c context.Context, ctx *app.RequestContext) {
	var user model.User

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.H{
			"resp": model.Response{
				Code: http.StatusBadGateway,
				Msg:  err.Error() + "参数错误",
				Data: "nil",
			},
		})
		return
	}

	rpcResp, err := global.Clients.UserClient.Register(c, &userrpc.RegisterReq{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{
			"resp": model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error() + "rpc服务错误",
				Data: "nil",
			},
		})
	}

	ctx.JSON(http.StatusOK, utils.H{
		"resp": model.Response{
			Code: 0,
			Msg:  "ok",
			Data: rpcResp.Resp.Data,
		},
	})
}

func Login(c context.Context, ctx *app.RequestContext) {
	var user model.User

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.H{
			"resp": model.Response{
				Code: http.StatusBadGateway,
				Msg:  err.Error() + "参数错误",
				Data: "nil",
			},
		})
		return
	}

	rpcResp, err := global.Clients.UserClient.Login(c, &userrpc.LoginReq{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{
			"resp": model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error() + "rpc服务错误",
				Data: "nil",
			},
		})
		return
	}

	token, err := my_utils.GenerateToken(rpcResp.Resp.Data)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"resp": model.Response{
			Code: 0,
			Msg:  "ok",
			Data: rpcResp.Resp.Data + "	" + token,
		},
	})

}
