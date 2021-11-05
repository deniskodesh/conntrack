package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {

			// conntrack_Total.Add(GetTableEntriesNumber())
			time.Sleep(3 * time.Second)

		}
	}()

	go func() {
		for {

			for ip, val := range HowMatches(GetRecordsFromTable()) {
				Top15.With(prometheus.Labels{"192.168.24.201": ip}).Set(float64(val))
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

var (
	conntrack_Total = promauto.NewCounter(prometheus.CounterOpts{
		Name: "conntrack_session_total",
		Help: "Shows current number records in table",
	})

	Top15 = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "job_session",
			Help: "Session info",
		},
		[]string{"192.168.24.201"},
	)
)

func main() {
	recordMetrics()

	println(GetTableEntriesNumber())

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
