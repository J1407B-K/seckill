package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"seckill/handler"
	middle "seckill/utils"
)

func InitRouter() {
	h := server.New(server.WithHostPorts(":8080"))
	h.PanicHandler = func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(500, utils.H{"msg": "panic"})
	}

	v1 := h.Group("/")
	{
		v1.POST("register", handler.Register)
		v1.POST("login", handler.Login)
	}

	v2 := h.Group("/")
	{
		v2.POST("createorder", handler.CreateOrder)
		v2.POST("confirmorder", handler.ConfirmOrder)
		v2.POST("queryorder", handler.QueryOrder)
	}

	v2.Use(middle.JWTAuthMiddleware())

	h.Spin()
}
