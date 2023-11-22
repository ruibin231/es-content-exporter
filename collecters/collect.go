package collecters

import (
	"es-content-export/pkgs/es"
	"es-content-export/settings"
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
	EsCollect.ESConn.WithLabelValues(settings.Config.Host).Set(0)
	EsCollect.LogQueryGauge.With(prometheus.Labels{
		"es": settings.Config.Host, "field": settings.Config.Field,
		"content": settings.Config.Content, "index": settings.Config.IndexPrefix}).Set(0)
	EsCollect.LogQueryTotal.With(prometheus.Labels{
		"es": settings.Config.Host, "field": settings.Config.Field,
		"content": settings.Config.Content, "index": settings.Config.IndexPrefix,
	}).Add(float64(0))
}

func TickerTask() {
	cycle := settings.Config.Cycle
	ticker := time.NewTicker(time.Duration(cycle) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		queryCount, err := es.QueryLogCount()
		if err != nil {
			EsCollect.ESConn.WithLabelValues(settings.Config.Host).Set(1)
		} else {
			EsCollect.ESConn.WithLabelValues(settings.Config.Host).Set(0)
		}
		EsCollect.LogQueryGauge.With(prometheus.Labels{
			"es": settings.Config.Host, "field": settings.Config.Field,
			"content": settings.Config.Content, "index": settings.Config.IndexPrefix}).
			Set(float64(queryCount))
		EsCollect.LogQueryTotal.With(prometheus.Labels{
			"es": settings.Config.Host, "field": settings.Config.Field,
			"content": settings.Config.Content, "index": settings.Config.IndexPrefix,
		}).Add(float64(queryCount))
	}
}

func init() {
	EsCollect = newMetrics()
}
