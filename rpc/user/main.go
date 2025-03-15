package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	user "seckill/idl/kitex_gen/user/userservice"
	"seckill/rpc/user/flag"
	"seckill/rpc/user/initialize"
)

func main() {
	initialize.SetupViper()
	db := initialize.InitGormDB()

	option := flag.Parse()
	ok := flag.DBOption(db, option)
	if !ok {
		log.Println("未自动建表")
	}

	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8890")

	svr := user.NewServer(&UserServiceImpl{DB: db},
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "userservice",
		}))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
