package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	fetchJobsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fetch_jobs_total",
			Help: "Total number of fetch jobs",
		},
		[]string{"status"},
	)

	productsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "products_total",
			Help: "Total number of products",
		},
	)

	productVersionsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "product_versions_total",
			Help: "Total number of product versions",
		},
		[]string{"product_id"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(fetchJobsTotal)
	prometheus.MustRegister(productsTotal)
	prometheus.MustRegister(productVersionsTotal)
}

func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(c.Writer.Status())),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}

func MetricsHandler() gin.HandlerFunc {
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func IncrementFetchJobMetric(status string) {
	fetchJobsTotal.WithLabelValues(status).Inc()
}

func SetProductsMetric(count float64) {
	productsTotal.Set(count)
}

func SetProductVersionsMetric(productID string, count float64) {
	productVersionsTotal.WithLabelValues(productID).Set(count)
}