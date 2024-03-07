package collecters

import (
	"es-content-export/pkgs/es"
	"es-content-export/settings"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var EsCollect *ElasticMetrics

type ElasticMetrics struct {
	ESConn        *prometheus.GaugeVec
	LogQueryGauge *prometheus.GaugeVec
	LogQueryTotal *prometheus.CounterVec
}

func newMetrics() *ElasticMetrics {
	return &ElasticMetrics{
		ESConn: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "es_query_failed",
			Help: "Current elastic search connection",
		},
			[]string{"node"}),
		LogQueryGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "es_log_alert_count",
			Help: "Query elastic search message errors count",
		},
			[]string{"es", "field", "content", "index"}),
		LogQueryTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "es_log_query_counter",
				Help: "Query elastic search message counter",
			},
			[]string{"es", "field", "content", "index"}),
	}
}

func RegistryEsCollect(reg *prometheus.Registry) {
	reg.MustRegister(EsCollect.ESConn, EsCollect.LogQueryTotal, EsCollect.LogQueryGauge)
}

func TickerTask(esClient *settings.EsClient, queryData *settings.QueryData) {
	cycle := queryData.Cycle
	ticker := time.NewTicker(time.Duration(cycle) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		queryCount, err := es.QueryLogCount(esClient, queryData)
		if err != nil {
			EsCollect.ESConn.WithLabelValues(esClient.Host).Set(1)
		} else {
			EsCollect.ESConn.WithLabelValues(esClient.Host).Set(0)
		}
		EsCollect.LogQueryGauge.With(prometheus.Labels{
			"es": esClient.Host, "field": queryData.Field,
			"content": queryData.Content, "index": queryData.IndexPrefix}).
			Set(float64(queryCount))
		EsCollect.LogQueryTotal.With(prometheus.Labels{
			"es": esClient.Host, "field": queryData.Field,
			"content": queryData.Content, "index": queryData.IndexPrefix,
		}).Add(float64(queryCount))
	}
}

func StartTasks() {
	for _, q := range settings.Config.QueryList {
		if e, ok := settings.ESMap[q.ES]; ok {
			go TickerTask(e, q)
		} else {
			fmt.Printf("找不到es: %s", q.ES)
		}
	}
}

func init() {
	EsCollect = newMetrics()
}
