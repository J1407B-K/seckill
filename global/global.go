package global

import (
	"github.com/cloudwego/kitex/pkg/discovery"
	"seckill/model"
)

var (
	Clients  = model.Clients{}
	Resolver *discovery.Resolver
)
