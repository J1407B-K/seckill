package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"seckill/handler"
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

	h.Spin()
}
