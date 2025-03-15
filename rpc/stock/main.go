package main

import (
	"log"
	stock "seckill/idl/kitex_gen/stock/stockservice"
)

func main() {
	svr := stock.NewServer(new(StockServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
