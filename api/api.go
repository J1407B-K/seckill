package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"seckill/handler"
	my_utils "seckill/utils"
)

func InitRouter() {
	h := server.New(server.WithHostPorts(":8080"))

	authMiddlewire, err := my_utils.NewMiddle()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	v1 := h.Group("/")
	{
		v1.POST("register", handler.Register)
		v1.POST("login", authMiddlewire.LoginHandler)
	}

	v2 := h.Group("/")
	v2.Use(authMiddlewire.MiddlewareFunc())
	{
		v2.POST("createorder", handler.CreateOrder)
		v2.POST("confirmorder", handler.ConfirmOrder)
		v2.POST("queryorder", handler.QueryOrder)
	}

	h.Spin()
}
