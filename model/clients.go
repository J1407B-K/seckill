package model

import (
	order "seckill/idl/kitex_gen/order/orderservice"
	stock "seckill/idl/kitex_gen/stock/stockservice"
	user "seckill/idl/kitex_gen/user/userservice"
)

type Clients struct {
	UserClient  user.Client
	StockClient stock.Client
	OrderClient order.Client
}
