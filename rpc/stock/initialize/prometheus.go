package initialize

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func InitPrometheus() *prometheus.CounterVec {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "stockservice",
		Help: "Total number of requests processed by the stock service",
	},
		[]string{"handler", "method"},
	)

	return requestCounter
}

func RegisterPromethus(requestCounter *prometheus.CounterVec) {
	prometheus.MustRegister(requestCounter)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Prometheus metrics server is running on :9102")
		if err := http.ListenAndServe(":9102", nil); err != nil {
			log.Fatalf("Failed to start Prometheus metrics server: %v", err)
		}
	}()
}
