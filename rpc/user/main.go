package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
	user "seckill/idl/kitex_gen/user/userservice"
	"seckill/rpc/user/init"
)

func main() {
	init.SetupViper()
	db := init.InitGormDB()

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
