package provider

import "github.com/prometheus/client_golang/prometheus"

// MetricsProvider defines the interface for collecting Starlink metrics
type MetricsProvider interface {
	// Collect returns metrics in Prometheus format
	Collect() error
	// GetRegistry returns the Prometheus registry with registered collectors
	GetRegistry() *prometheus.Registry
	// Close performs any necessary cleanup
	Close() error
}
