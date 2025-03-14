package main

import (
	"log"
	order "seckill/service/order_service/kitex_gen/order/orderservice"
)

func main() {
	svr := order.NewServer(new(OrderServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
