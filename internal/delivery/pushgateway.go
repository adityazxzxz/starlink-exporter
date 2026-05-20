package delivery

import (
	"fmt"
	"time"

	"github.com/danopstech/starlink_exporter/internal/provider"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
)

// PushgatewayDelivery pushes metrics to Prometheus Pushgateway
type PushgatewayDelivery struct {
	config   *Config
	ticker   *time.Ticker
	done     chan bool
	interval time.Duration
}

// NewPushgatewayDelivery creates a new pushgateway delivery mode
func NewPushgatewayDelivery(config *Config) (*PushgatewayDelivery, error) {
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return nil, fmt.Errorf("invalid interval %s: %v", config.Interval, err)
	}

	return &PushgatewayDelivery{
		config:   config,
		interval: interval,
		done:     make(chan bool, 1),
	}, nil
}

// Run starts pushing metrics to pushgateway periodically
func (p *PushgatewayDelivery) Run(metricsProvider provider.MetricsProvider) error {
	log.Infof("starting pushgateway delivery mode: %s (job=%s, instance=%s, interval=%s)",
		p.config.PushgatewayURL, p.config.Job, p.config.Instance, p.config.Interval)

	// Initial push immediately
	if err := p.push(metricsProvider); err != nil {
		log.Warnf("initial push failed: %v", err)
	} else {
		log.Infof("initial metrics push successful")
	}

	// Create ticker for periodic pushes
	p.ticker = time.NewTicker(p.interval)

	// Run push loop
	for {
		select {
		case <-p.ticker.C:
			if err := p.push(metricsProvider); err != nil {
				log.Errorf("push to pushgateway failed: %v", err)
			} else {
				log.Infof("metrics pushed to pushgateway successfully")
			}
		case <-p.done:
			log.Info("shutting down pushgateway delivery")
			// Final push before shutdown
			if err := p.push(metricsProvider); err != nil {
				log.Warnf("final push on shutdown failed: %v", err)
			}
			return nil
		}
	}
}

// push sends metrics to pushgateway
func (p *PushgatewayDelivery) push(metricsProvider provider.MetricsProvider) error {
	pusher := push.New(p.config.PushgatewayURL, p.config.Job).
		Gatherer(metricsProvider.GetRegistry()).
		Grouping("instance", p.config.Instance)

	return pusher.Push()
}

// Shutdown gracefully shuts down pushgateway delivery
func (p *PushgatewayDelivery) Shutdown() error {
	if p.ticker != nil {
		p.ticker.Stop()
	}

	select {
	case p.done <- true:
	default:
	}

	return nil
}
