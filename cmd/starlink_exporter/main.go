package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/danopstech/starlink_exporter/internal/delivery"
	"github.com/danopstech/starlink_exporter/internal/exporter"
	"github.com/danopstech/starlink_exporter/internal/provider"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Parse command line arguments
	mode := flag.String("mode", "web", "delivery mode: web or pushgateway")
	source := flag.String("source", "live", "metrics source: live or dummy")
	listen := flag.String("listen", ":9817", "listen address for web mode (e.g., :9817)")
	address := flag.String("address", exporter.DishAddress, "IP address and port to reach dish (live mode only)")
	pushgatewayURL := flag.String("pushgateway", "", "Pushgateway URL (e.g., http://pushgateway:9091)")
	job := flag.String("job", "starlink_exporter", "job name for pushgateway")
	instance := flag.String("instance", "starlink_dish", "instance name for pushgateway")
	interval := flag.String("interval", "15s", "push interval for pushgateway mode")
	logLevel := flag.String("log-level", "info", "log level: debug, info, warn, error")
	flag.Parse()

	// Set log level
	if lvl, err := log.ParseLevel(*logLevel); err == nil {
		log.SetLevel(lvl)
	}

	// Validate mode and source
	if *mode != "web" && *mode != "pushgateway" {
		log.Fatalf("invalid mode: %s (valid: web, pushgateway)", *mode)
	}
	if *source != "live" && *source != "dummy" {
		log.Fatalf("invalid source: %s (valid: live, dummy)", *source)
	}

	log.Infof("starting starlink exporter - mode=%s, source=%s", *mode, *source)

	// Create metrics provider based on source
	var metricsProvider provider.MetricsProvider
	var err error

	if *source == "live" {
		metricsProvider, err = provider.NewLiveMetricsProvider(*address)
		if err != nil {
			log.Fatalf("failed to create live metrics provider: %v", err)
		}
	} else {
		metricsProvider = provider.NewDummyMetricsProvider()
	}
	defer func() {
		if err := metricsProvider.Close(); err != nil {
			log.Errorf("error closing metrics provider: %v", err)
		}
	}()

	// Create delivery mode based on mode
	var deliveryMode delivery.DeliveryMode

	config := &delivery.Config{
		Listen:         *listen,
		PushgatewayURL: *pushgatewayURL,
		Job:            *job,
		Instance:       *instance,
		Interval:       *interval,
	}

	if *mode == "web" {
		deliveryMode = delivery.NewWebDelivery(config)
	} else {
		deliveryMode, err = delivery.NewPushgatewayDelivery(config)
		if err != nil {
			log.Fatalf("failed to create pushgateway delivery: %v", err)
		}
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run delivery mode in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- deliveryMode.Run(metricsProvider)
	}()

	// Wait for shutdown signal
	select {
	case sig := <-sigChan:
		log.Infof("received signal: %v", sig)
		if err := deliveryMode.Shutdown(); err != nil {
			log.Errorf("error during shutdown: %v", err)
		}
	case err := <-errChan:
		if err != nil {
			log.Fatalf("delivery mode error: %v", err)
		}
	}

	log.Info("exporter stopped")
}
