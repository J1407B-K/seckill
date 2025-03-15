package initialize

import (
	"github.com/cloudwego/kitex/client"
	"seckill/global"
	"seckill/idl/kitex_gen/order/orderservice"
	"seckill/idl/kitex_gen/stock/stockservice"
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

	sc, err := stockservice.NewClient("stockservice", client.WithResolver(*global.Resolver))
	if err != nil {
		panic(err)
	}
	global.Clients.StockClient = sc

	return nil
}
