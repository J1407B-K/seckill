package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	stock "seckill/idl/kitex_gen/stock/stockservice"
	"seckill/rpc/stock/flag"
	"seckill/rpc/stock/initialize"
)

func main() {
	initialize.SetupViper()
	db := initialize.InitGormDB()
	rdb := initialize.InitRedisDB()
	redsync := initialize.InitRedisSync()

	option := flag.Parse()
	ok := flag.DBOption(db, option)
	if !ok {
		log.Println("未自动建表")
	}

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8891")
	svr := stock.NewServer(&StockServiceImpl{db: db, rdb: rdb, redsync: redsync},
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "stockservice",
		}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
