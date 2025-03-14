package main

import (
	"log"
	stock "seckill/service/stock_service/kitex_gen/stock/stockservice"
)

func main() {
	svr := stock.NewServer(new(StockServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
