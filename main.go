package main

import (
	"container/heap"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			_, sessionsCount := GetRecordsFromTable()
			conntrack_Total.Add(sessionsCount)
			time.Sleep(3 * time.Second)
			println(sessionsCount)
		}
	}()

	go func() {
		for {
			sessions, _ := GetRecordsFromTable()

			// Create a heap from the map and print the top N values.
			h := getHeap(HowMatches(sessions))
			n := 15
			for i := 0; i < n; i++ {
				fmt.Printf("%d) %#v\n", i+1, heap.Pop(h))
			}

			for ip, val := range HowMatches(sessions) {
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


	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
