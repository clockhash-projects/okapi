package logger

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"okapi/internal/models"
)

var (
	// ServiceStatus indicates the current status of a service (0=unknown, 1=operational, 2=degraded, 3=partial_outage, 4=major_outage)
	ServiceStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "okapi_service_status",
		Help: "Current status of the service (1=operational, 2=degraded, 3=partial_outage, 4=major_outage)",
	}, []string{"service"})

	// PollDuration tracks how long it takes to fetch status from an adapter
	PollDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "okapi_poll_duration_seconds",
		Help:    "Time taken to poll service status",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "success"})
)

// RecordStatus updates the status gauge for a service
func RecordStatus(service string, status models.Status) {
	val := 0.0
	switch status {
	case models.StatusOperational:
		val = 1.0
	case models.StatusDegraded:
		val = 2.0
	case models.StatusPartialOutage:
		val = 3.0
	case models.StatusMajorOutage:
		val = 4.0
	default:
		val = 0.0
	}
	ServiceStatus.WithLabelValues(service).Set(val)
}
