package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	order "seckill/idl/kitex_gen/order/orderservice"
	"seckill/rpc/order/flag"
	"seckill/rpc/order/initialize"
)

func main() {
	initialize.SetupViper()
	db := initialize.InitGormDB()

	option := flag.Parse()
	ok := flag.DBOption(db, option)
	if !ok {
		log.Println("未自动建表")
	}

	initialize.InitTracer()

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	stockCli, err := NewStockClient()
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8892")
	svr := order.NewServer(&OrderServiceImpl{db: db, stockCli: stockCli},
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "orderservice",
		}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
