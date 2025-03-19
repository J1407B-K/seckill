package global

import (
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/prometheus/client_golang/prometheus"
	"seckill/model"
)

var (
	Clients  = model.Clients{}
	Resolver *discovery.Resolver
	Cv       *prometheus.CounterVec
)
