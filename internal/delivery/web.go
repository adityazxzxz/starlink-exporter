package delivery

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danopstech/starlink_exporter/internal/provider"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// WebDelivery exposes metrics via HTTP /metrics endpoint
type WebDelivery struct {
	config *Config
	server *http.Server
}

// NewWebDelivery creates a new web delivery mode
func NewWebDelivery(config *Config) *WebDelivery {
	return &WebDelivery{
		config: config,
	}
}

// Run starts the HTTP server and exposes metrics
func (w *WebDelivery) Run(metricsProvider provider.MetricsProvider) error {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(writer, "OK\n")
	})

	// Root page
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`<html>
		 <head><title>Starlink Exporter</title></head>
		 <body>
		 <h1>Starlink Exporter</h1>
		 <p><a href='/metrics'>Metrics</a></p>
		 <p><a href='/health'>Health</a></p>
		 </body>
		 </html>`))
	})

	// Metrics endpoint
	mux.Handle("/metrics", promhttp.HandlerFor(
		metricsProvider.GetRegistry(),
		promhttp.HandlerOpts{},
	))

	w.server = &http.Server{
		Addr:    w.config.Listen,
		Handler: mux,
	}

	log.Infof("starting web delivery mode on %s", w.config.Listen)

	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("web server error: %v", err)
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (w *WebDelivery) Shutdown() error {
	if w.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return w.server.Shutdown(ctx)
}
