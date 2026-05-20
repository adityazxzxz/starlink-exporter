package provider

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// DummyMetricsProvider generates realistic dummy Starlink metrics for testing
type DummyMetricsProvider struct {
	registry *prometheus.Registry
	mu       sync.RWMutex

	// Metrics
	dishUpMetric                 prometheus.Gauge
	dishScrapeDurationSeconds    prometheus.Gauge
	dishUptimeSeconds            prometheus.Gauge
	dishCellId                   prometheus.Gauge
	dishPopPingDropRatio         prometheus.Gauge
	dishPopPingLatencySeconds    prometheus.Gauge
	dishSnr                      prometheus.Gauge
	dishUplinkThroughputBytes    prometheus.Gauge
	dishDownlinkThroughputBytes  prometheus.Gauge
	dishCurrentlyObstructed      prometheus.Gauge
	dishFractionObstructionRatio prometheus.Gauge
	dishLast24hObstructedSeconds prometheus.Gauge
	dishState                    prometheus.Gauge
	dishAlertMotorsStuck         prometheus.Gauge
	dishAlertThermalThrottle     prometheus.Gauge
	dishAlertThermalShutdown     prometheus.Gauge
	dishAlertMastNotNearVertical prometheus.Gauge
	dishUnexpectedLocation       prometheus.Gauge
	dishSlowEthernetSpeeds       prometheus.Gauge

	// Info metric
	dishInfoMetric prometheus.GaugeVec

	// Ticker for periodic updates
	ticker *time.Ticker
	done   chan bool

	// State tracking
	startTime time.Time
}

// NewDummyMetricsProvider creates a new dummy metrics provider
func NewDummyMetricsProvider() *DummyMetricsProvider {
	registry := prometheus.NewRegistry()

	provider := &DummyMetricsProvider{
		registry:  registry,
		ticker:    time.NewTicker(5 * time.Second),
		done:      make(chan bool, 1),
		startTime: time.Now(),
	}

	// Initialize all metrics
	provider.dishUpMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "up",
		Help:      "Was the last query of Starlink dish successful.",
	})

	provider.dishScrapeDurationSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "scrape_duration_seconds",
		Help:      "Time to scrape metrics from starlink dish",
	})

	provider.dishUptimeSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "uptime_seconds",
		Help:      "Dish running time",
	})

	provider.dishCellId = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "cell_id",
		Help:      "Cell ID dish is located in",
	})

	provider.dishPopPingDropRatio = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "pop_ping_drop_ratio",
		Help:      "Percent of pings dropped",
	})

	provider.dishPopPingLatencySeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "pop_ping_latency_seconds",
		Help:      "Latency of connection in seconds",
	})

	provider.dishSnr = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "snr",
		Help:      "Signal strength of the connection",
	})

	provider.dishUplinkThroughputBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "uplink_throughput_bytes",
		Help:      "Amount of bandwidth in bytes per second upload",
	})

	provider.dishDownlinkThroughputBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "downlink_throughput_bytes",
		Help:      "Amount of bandwidth in bytes per second download",
	})

	provider.dishCurrentlyObstructed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "currently_obstructed",
		Help:      "Status of view of the sky",
	})

	provider.dishFractionObstructionRatio = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "fraction_obstruction_ratio",
		Help:      "Percentage of obstruction",
	})

	provider.dishLast24hObstructedSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "last_24h_obstructed_seconds",
		Help:      "Number of seconds view of sky has been obstructed in the last 24hours",
	})

	provider.dishState = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "state",
		Help:      "The current state of the Dish (Unknown, Booting, Searching, Connected).",
	})

	provider.dishAlertMotorsStuck = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "alert_motors_stuck",
		Help:      "Status of motor stuck",
	})

	provider.dishAlertThermalThrottle = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "alert_thermal_throttle",
		Help:      "Status of thermal throttling",
	})

	provider.dishAlertThermalShutdown = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "alert_thermal_shutdown",
		Help:      "Status of thermal shutdown",
	})

	provider.dishAlertMastNotNearVertical = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "alert_mast_not_near_vertical",
		Help:      "Status of mast position",
	})

	provider.dishUnexpectedLocation = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "alert_unexpected_location",
		Help:      "Status of location",
	})

	provider.dishSlowEthernetSpeeds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "alert_slow_eth_speeds",
		Help:      "Status of ethernet",
	})

	provider.dishInfoMetric = *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "starlink",
		Subsystem: "dish",
		Name:      "info",
		Help:      "Running software versions and IDs of hardware",
	}, []string{"device_id", "hardware_version", "software_version", "country_code", "utc_offset"})

	// Register all metrics
	registry.MustRegister(provider.dishUpMetric)
	registry.MustRegister(provider.dishScrapeDurationSeconds)
	registry.MustRegister(provider.dishUptimeSeconds)
	registry.MustRegister(provider.dishCellId)
	registry.MustRegister(provider.dishPopPingDropRatio)
	registry.MustRegister(provider.dishPopPingLatencySeconds)
	registry.MustRegister(provider.dishSnr)
	registry.MustRegister(provider.dishUplinkThroughputBytes)
	registry.MustRegister(provider.dishDownlinkThroughputBytes)
	registry.MustRegister(provider.dishCurrentlyObstructed)
	registry.MustRegister(provider.dishFractionObstructionRatio)
	registry.MustRegister(provider.dishLast24hObstructedSeconds)
	registry.MustRegister(provider.dishState)
	registry.MustRegister(provider.dishAlertMotorsStuck)
	registry.MustRegister(provider.dishAlertThermalThrottle)
	registry.MustRegister(provider.dishAlertThermalShutdown)
	registry.MustRegister(provider.dishAlertMastNotNearVertical)
	registry.MustRegister(provider.dishUnexpectedLocation)
	registry.MustRegister(provider.dishSlowEthernetSpeeds)
	registry.MustRegister(&provider.dishInfoMetric)

	// Set initial values and start update loop
	provider.updateMetrics()
	go provider.startPeriodicUpdates()

	log.Info("dummy metrics provider initialized")

	return provider
}

// startPeriodicUpdates updates metrics periodically
func (p *DummyMetricsProvider) startPeriodicUpdates() {
	for {
		select {
		case <-p.ticker.C:
			p.updateMetrics()
		case <-p.done:
			return
		}
	}
}

// updateMetrics updates all dummy metrics with realistic values
func (p *DummyMetricsProvider) updateMetrics() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Always up
	p.dishUpMetric.Set(1.0)

	// Scrape duration
	p.dishScrapeDurationSeconds.Set(float64(rand.Intn(50)+10) / 1000.0)

	// Uptime: increases over time with random variations
	uptime := time.Since(p.startTime).Seconds() + float64(rand.Intn(1000)-500)
	p.dishUptimeSeconds.Set(math.Max(0, uptime))

	// Cell ID
	p.dishCellId.Set(float64(rand.Intn(50000) + 1000))

	// Ping drop ratio (0-5%)
	p.dishPopPingDropRatio.Set(float64(rand.Intn(500)) / 10000.0)

	// Latency (20-120 ms)
	latency := float64(rand.Intn(100)+20) / 1000.0
	p.dishPopPingLatencySeconds.Set(latency)

	// SNR (5-15 dB)
	snr := 5.0 + rand.Float64()*10.0
	p.dishSnr.Set(snr)

	// Upload throughput (5-40 Mbps converted to bytes)
	uploadMbps := float64(rand.Intn(35) + 5)
	p.dishUplinkThroughputBytes.Set(uploadMbps * 125000) // Mbps to bytes/sec

	// Download throughput (50-250 Mbps converted to bytes)
	downloadMbps := float64(rand.Intn(200) + 50)
	p.dishDownlinkThroughputBytes.Set(downloadMbps * 125000) // Mbps to bytes/sec

	// Currently obstructed (mostly 0, occasionally 1)
	obstructed := float64(0)
	if rand.Intn(100) < 5 {
		obstructed = 1.0
	}
	p.dishCurrentlyObstructed.Set(obstructed)

	// Obstruction ratio (0-20%)
	p.dishFractionObstructionRatio.Set(float64(rand.Intn(2000)) / 10000.0)

	// Last 24h obstructed seconds (0-3600 seconds = 0-1 hour)
	p.dishLast24hObstructedSeconds.Set(float64(rand.Intn(3600)))

	// Dish state (3 = connected)
	p.dishState.Set(3.0)

	// Alerts (all false/0)
	p.dishAlertMotorsStuck.Set(0.0)
	p.dishAlertThermalThrottle.Set(0.0)
	p.dishAlertThermalShutdown.Set(0.0)
	p.dishAlertMastNotNearVertical.Set(0.0)
	p.dishUnexpectedLocation.Set(0.0)
	p.dishSlowEthernetSpeeds.Set(0.0)

	// Set info metric
	p.dishInfoMetric.WithLabelValues(
		"dummy-device-001",
		"rev2_proto_v2",
		"2024.07.01",
		"US",
		"-8",
	).Set(1.0)
}

// Collect performs the actual collection (no-op for dummy provider)
func (p *DummyMetricsProvider) Collect() error {
	// Metrics are updated periodically by the ticker
	return nil
}

// GetRegistry returns the Prometheus registry
func (p *DummyMetricsProvider) GetRegistry() *prometheus.Registry {
	return p.registry
}

// Close stops the periodic updates
func (p *DummyMetricsProvider) Close() error {
	p.ticker.Stop()
	select {
	case p.done <- true:
	default:
	}
	return nil
}
