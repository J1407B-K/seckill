package initialize

import (
	"github.com/cloudwego/kitex/client"
	"seckill/global"
	"seckill/idl/kitex_gen/order/orderservice"
	"seckill/idl/kitex_gen/user/userservice"
)

func InitNewClient() error {
	uc, err := userservice.NewClient("userservice", client.WithResolver(*global.Resolver))
	if err != nil {
		panic(err)
	}
	global.Clients.UserClient = uc

	oc, err := orderservice.NewClient("orderservice", client.WithResolver(*global.Resolver))
	if err != nil {
		panic(err)
	}
	global.Clients.OrderClient = oc
	return nil
}
