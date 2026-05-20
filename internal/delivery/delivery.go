package delivery

import "github.com/danopstech/starlink_exporter/internal/provider"

// DeliveryMode defines the interface for delivering metrics
type DeliveryMode interface {
	// Run starts the delivery mode and blocks until shutdown
	Run(metricsProvider provider.MetricsProvider) error
	// Shutdown gracefully shuts down the delivery mode
	Shutdown() error
}

// Config contains common delivery mode configuration
type Config struct {
	// Listen address for web server (e.g., ":9817")
	Listen string
	// Pushgateway URL (e.g., "http://pushgateway:9091")
	PushgatewayURL string
	// Job name for pushgateway
	Job string
	// Instance name for pushgateway
	Instance string
	// Interval for pushing metrics
	Interval string
}
