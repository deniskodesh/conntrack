package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	settings Config

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

	Top = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "all_ips_sessions",
			Help: "Show all IPs with session count",
		},
		[]string{"ip"},
	)
)

func init() {

	log.WithFields(log.Fields{}).Info("Starting initializing")
	//enable JSON format
	log.SetFormatter(&log.JSONFormatter{})

	conf := readFromFile("./settings.json")
	log.WithFields(log.Fields{}).Info("Reading configuration file")
	json.Unmarshal([]byte(conf), &settings)

	if settings.LogDebug {
		//Shows caller
		log.SetReportCaller(true)
	}

	// Metrics have to be registered to be exposed:
	log.WithFields(log.Fields{}).Info("Registering prometheus metrics")
	prometheus.MustRegister(conntrack_Total)
	prometheus.MustRegister(conntrack_Max)
	prometheus.MustRegister(Top)

}

func main() {
	log.Info("Getting prometheus metrics")
	recordMetrics()

	http.Handle(settings.MetricsRoutePath, promhttp.Handler())

	log.WithFields(log.Fields{
		"port": settings.Port,
		"path": settings.MetricsRoutePath,
	}).Info("Start http server")
	http.ListenAndServe(":"+settings.Port, nil)
}
