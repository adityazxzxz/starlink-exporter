package provider

import (
	"github.com/danopstech/starlink_exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// LiveMetricsProvider collects metrics from real Starlink device
type LiveMetricsProvider struct {
	exporter *exporter.Exporter
	registry *prometheus.Registry
}

// NewLiveMetricsProvider creates a new live metrics provider that connects to Starlink dish
func NewLiveMetricsProvider(address string) (*LiveMetricsProvider, error) {
	exp, err := exporter.New(address)
	if err != nil {
		log.Errorf("failed to create exporter: %v", err)
		return nil, err
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(exp)

	log.Infof("connected to starlink dish at %s with ID: %s", address, exp.DishID)

	return &LiveMetricsProvider{
		exporter: exp,
		registry: registry,
	}, nil
}

// Collect performs the actual collection from the dish
func (p *LiveMetricsProvider) Collect() error {
	// The actual collection happens through the prometheus registry
	// when metrics are scraped or pushed
	return nil
}

// GetRegistry returns the Prometheus registry
func (p *LiveMetricsProvider) GetRegistry() *prometheus.Registry {
	return p.registry
}

// Close closes the gRPC connection
func (p *LiveMetricsProvider) Close() error {
	if p.exporter != nil && p.exporter.Conn != nil {
		return p.exporter.Conn.Close()
	}
	return nil
}
