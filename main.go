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
			conntrack_Total.Add(StringToFloat(string(readFromFile("./1.txt"))))
			time.Sleep(3 * time.Second)

		}
	}()

	go func() {
		for {

			for ip, val := range GetRecordsFromTable() {
				Top15.With(prometheus.Labels{"IP": ip}).Set(float64(val))
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
			Name: "job_location",
			Help: "Location connections number",
		},
		[]string{"location"},
	)
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}