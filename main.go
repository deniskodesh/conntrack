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
			conntrack_Total.Set(Float64frombytes(fileBytes))

			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {

			fileBytes := readFromFile("/proc/sys/net/netfilter/nf_conntrack_max")
			conntrack_Max.Set(Float64frombytes(fileBytes))

			time.Sleep(3 * time.Second)
		}
	}()
}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(conntrack_Total)
	prometheus.MustRegister(conntrack_Max)

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

	conntrack_Max = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "conntrack_session_max",
			Help: "Shows max number records in table from file /proc/sys/net/netfilter/nf_conntrack_max",
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
