package initialize

import (
	"github.com/cloudwego/kitex/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.0"
	"log"
	"net/http"
	"seckill/global"
	"seckill/idl/kitex_gen/order/orderservice"
	"seckill/idl/kitex_gen/user/userservice"
)

func InitNewClient() error {
	uc, err := userservice.NewClient("userservice", client.WithResolver(*global.Resolver))
	if err != nil {
		panic(err)
	}
	global.Clients.UserClient = uc

	oc, err := orderservice.NewClient("orderservice", client.WithResolver(*global.Resolver))
	if err != nil {
		panic(err)
	}
	global.Clients.OrderClient = oc
	return nil
}

func InitTracer() {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/traces")))
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("client"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func InitPrometheus() *prometheus.CounterVec {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ApiRequest",
		Help: "Total number of requests processed by the api service",
	},
		[]string{"handler", "method"},
	)

	return requestCounter
}

func RegisterPromethus(requestCounter *prometheus.CounterVec) {
	prometheus.MustRegister(requestCounter)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Prometheus metrics server is running on :9100")
		if err := http.ListenAndServe(":9100", nil); err != nil {
			log.Fatalf("Failed to start Prometheus metrics server: %v", err)
		}
	}()
}
