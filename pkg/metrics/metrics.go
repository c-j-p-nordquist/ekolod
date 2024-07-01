package metrics

import (
	"net/http"

	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/proberesult"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"target", "path", "method"})

	HttpResponseSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_response_size_bytes",
		Help: "Size of HTTP responses in bytes.",
	}, []string{"target", "path"})

	TLSVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tls_version",
		Help: "TLS version used for the connection.",
	}, []string{"target", "version"})

	CertExpiryDays = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cert_expiry_days",
		Help: "Number of days until the SSL certificate expires.",
	}, []string{"target"})
)

func InitMetrics() {
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpResponseSize)
	prometheus.MustRegister(TLSVersion)
	prometheus.MustRegister(CertExpiryDays)
}

func UpdatePrometheusMetrics(target *config.Target, check config.Check, result *proberesult.ProbeResult) {
	HttpRequestDuration.With(prometheus.Labels{
		"target": target.Name,
		"path":   check.Path,
		"method": "GET",
	}).Observe(result.Duration)

	HttpResponseSize.With(prometheus.Labels{
		"target": target.Name,
		"path":   check.Path,
	}).Set(float64(result.ContentLength))

	if result.TLSVersion != "" {
		TLSVersion.With(prometheus.Labels{
			"target":  target.Name,
			"version": result.TLSVersion,
		}).Set(1)
	}

	if result.CertExpiryDays != 0 {
		CertExpiryDays.With(prometheus.Labels{
			"target": target.Name,
		}).Set(float64(result.CertExpiryDays))
	}
}

func Handler() http.Handler {
	return promhttp.Handler()
}
