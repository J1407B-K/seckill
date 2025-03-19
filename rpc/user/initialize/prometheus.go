package initialize

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func InitPrometheus() *prometheus.CounterVec {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "userservice",
		Help: "Total number of requests processed by the user service",
	},
		[]string{"handler", "method"},
	)

	return requestCounter
}

func RegisterPromethus(requestCounter *prometheus.CounterVec) {
	prometheus.MustRegister(requestCounter)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Prometheus metrics server is running on :9101")
		if err := http.ListenAndServe(":9101", nil); err != nil {
			log.Fatalf("Failed to start Prometheus metrics server: %v", err)
		}
	}()
}
