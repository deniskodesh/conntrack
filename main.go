package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {

			fileBytes := readFromFile("/proc/sys/net/netfilter/nf_conntrack_count")
			conntrack_Total.Set(float64(byteToInt(fileBytes)))

			time.Sleep(3 * time.Second)
		}
	}()

	// go func() {
	// 	for {
	// 		sessions, _ := GetRecordsFromTable()

	// 		getTopValues(5, sessions)
	// 		for ip, val := range HowMatches(sessions) {
	// 			Top15.With(prometheus.Labels{"192.168.24.201": ip}).Set(float64(val))

	// 		}
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()
}

var (
	// conntrack_Total = promauto.NewCounter(prometheus.CounterOpts{
	// 	Name: "conntrack_session_total",
	// 	Help: "Shows current number records in table",
	// })

	conntrack_Total = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "conntrack_session_total",
			Help: "Shows current number records in table from file /proc/sys/net/netfilter/nf_conntrack_count",
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

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
