package main

import (
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"seckill/api"
	"seckill/global"
	"seckill/initialize"
)

func main() {
	resolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	global.Resolver = &resolver

	err = initialize.InitNewClient()
	if err != nil {
		panic(err)
	}

	initialize.InitTracer()

	global.Cv = initialize.InitPrometheus()
	initialize.RegisterPromethus(global.Cv)

	api.InitRouter()
}
