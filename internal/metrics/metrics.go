package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"target", "path", "method"})
)

func InitMetrics() {
	prometheus.MustRegister(HttpRequestDuration)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
