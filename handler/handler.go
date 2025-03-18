package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"seckill/global"
	"seckill/idl/kitex_gen/order"
	userrpc "seckill/idl/kitex_gen/user"
	"seckill/model"
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

func Login(c context.Context, ctx *app.RequestContext) (interface{}, error) {
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
		return nil, nil
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
		return nil, nil
	}

	ctx.JSON(http.StatusOK, utils.H{
		"resp": model.Response{
			Code: 0,
			Msg:  "ok",
			Data: rpcResp.Resp.Data,
		},
	})
	return rpcResp.Resp.Data, nil
}

func CreateOrder(c context.Context, ctx *app.RequestContext) {
	var co model.CreateOrder

	user, ok := ctx.Get("userid")
	userinfo := user.(*model.User)
	err := ctx.BindJSON(&co)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.H{
			"resp": model.Response{
				Code: http.StatusBadGateway,
				Msg:  err.Error(),
				Data: "nil",
			},
		})
		return
	}

	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.H{
			"resp": model.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Get user error",
				Data: "nil",
			},
		})
		return
	}

	createOrderResp, err := global.Clients.OrderClient.CreateOrder(c, &order.OrderReq{
		UserId:    userinfo.UserId,
		ProductId: co.ProductId,
		Count:     int32(co.Count),
	})
	if err != nil {
		return
	}

	fmt.Println(createOrderResp)

	ctx.JSON(http.StatusOK, utils.H{
		"resp": model.Response{
			Code: 0,
			Msg:  "ok",
			Data: *createOrderResp.OrderId,
		},
	})
	return
}

func ConfirmOrder(c context.Context, ctx *app.RequestContext) {
	var co model.CofirmOrder
	err := ctx.BindJSON(&co)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.H{
			"resp": model.Response{
				Code: http.StatusBadGateway,
				Msg:  err.Error(),
				Data: "nil",
			},
		})
		return
	}

	confirmOrderResp, err := global.Clients.OrderClient.ConfirmOrder(c, co.OrderId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{
			"resp": model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: "nil",
			},
		})
		return
	}

	fmt.Println(confirmOrderResp)
	ctx.JSON(http.StatusOK, utils.H{
		"resp": model.Response{
			Code: 0,
			Msg:  "ok",
			Data: *confirmOrderResp.OrderId,
		},
	})
}

func QueryOrder(c context.Context, ctx *app.RequestContext) {
	user, ok := ctx.Get("userid")
	userinfo := user.(*model.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.H{
			"resp": model.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Get user error",
				Data: "nil",
			},
		})
		return
	}

	qo, err := global.Clients.OrderClient.QueryOrder(c, userinfo.UserId)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"resp": model.Response{
			Code: 0,
			Msg:  "ok",
			Data: qo.Message,
		},
	})
}
