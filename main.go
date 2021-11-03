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
			conntrack_Total.Add(StringToFloat(string(readFromFile("/proc/sys/net/netfilter/nf_conntrack_count"))))
			time.Sleep(3 * time.Second)

		}
	}()

	go func() {
		for {
			conntrack_Max.Add(StringToFloat(string(readFromFile("/proc/sys/net/netfilter/nf_conntrack_max"))))
			time.Sleep(3 * time.Second)

		}
	}()

	go func() {
		for {

			for ip, val := range HowMatches(GetRecordsFromTable()) {
				Top.With(prometheus.Labels{"ip": ip}).Set(float64(val))
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

	conntrack_Max = promauto.NewCounter(prometheus.CounterOpts{
		Name: "conntrack_session_max",
		Help: "Shows limits on OS",
	})

	Top = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "job_session",
			Help: "Session info",
		},
		[]string{"ip"},
	)
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
