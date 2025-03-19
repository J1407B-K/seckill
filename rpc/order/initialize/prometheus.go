package initialize

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func InitPrometheus() *prometheus.CounterVec {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "orderservice",
		Help: "Total number of requests processed by the order service",
	},
		[]string{"handler", "method"},
	)

	return requestCounter
}

func RegisterPromethus(requestCounter *prometheus.CounterVec) {
	prometheus.MustRegister(requestCounter)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Prometheus metrics server is running on :9103")
		if err := http.ListenAndServe(":9103", nil); err != nil {
			log.Fatalf("Failed to start Prometheus metrics server: %v", err)
		}
	}()
}
