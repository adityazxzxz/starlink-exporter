<p align="center">
  <img alt="logo" src="https://github.com/danopstech/starlink_exporter/raw/main/.docs/assets/logo.jpg" height="150" />
  <h3 align="center">Starlink Prometheus Exporter - Enhanced Edition</h3>
</p>

---

A flexible [Starlink](https://www.starlink.com/) exporter for Prometheus with support for multiple delivery modes and data sources. Not affiliated with or acting on behalf of Starlink(™)

[![goreleaser](https://github.com/danopstech/starlink_exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/danopstech/starlink_exporter/actions/workflows/release.yaml)
[![build](https://github.com/danopstech/starlink_exporter/actions/workflows/build.yaml/badge.svg)](https://github.com/danopstech/starlink_exporter/actions/workflows/build.yaml)
[![License](https://img.shields.io/github/license/danopstech/starlink_exporter)](/LICENSE)
[![Release](https://img.shields.io/github/release/danopstech/starlink_exporter.svg)](https://github.com/danopstech/starlink_exporter/releases/latest)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/danopstech/starlink_exporter)
![os/arch](https://img.shields.io/badge/os%2Farch-amd64-yellow)
![os/arch](https://img.shields.io/badge/os%2Farch-arm64-yellow)
![os/arch](https://img.shields.io/badge/os%2Farch-armv7-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/danopstech/starlink_exporter)](https://goreportcard.com/report/github.com/danopstech/starlink_exporter)

## Features

- **Flexible Delivery Modes**:
  - **Web Mode**: Expose metrics via HTTP `/metrics` endpoint for direct Prometheus scraping
  - **Pushgateway Mode**: Push metrics to Prometheus Pushgateway at configurable intervals
  
- **Multiple Data Sources**:
  - **Live**: Real metrics from Starlink dish via gRPC
  - **Dummy**: Realistic simulated metrics for testing and development
  
- **Production Ready**:
  - Graceful shutdown handling (SIGTERM, SIGINT)
  - Thread-safe metric generation
  - Cross-platform build support (AMD64, ARM64, ARMv7)
  - Docker support with multi-stage builds
  - Systemd service files included
  
- **Backward Compatible**:
  - Existing metric names unchanged
  - Existing Grafana dashboards continue to work
  - Same API contract

## Architecture

The exporter uses a clean, modular architecture:

```
┌─────────────────────┐
│  Metrics Provider   │
├─────────────────────┤
│ • LiveProvider      │  (connects to real dish)
│ • DummyProvider     │  (generates test metrics)
└─────────────────────┘
          ▼
┌─────────────────────┐
│  Delivery Mode      │
├─────────────────────┤
│ • WebDelivery       │  (HTTP /metrics endpoint)
│ • PushgatewayDel.   │  (push to pushgateway)
└─────────────────────┘
```

## Quick Start

### 1. Build from Source

```bash
# Clone the repository
git clone https://github.com/danopstech/starlink_exporter.git
cd starlink_exporter

# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific ARM targets (Raspberry Pi)
make build-arm64
make build-armv7
```

### 2. Web Mode (Direct Scraping)

#### Live Metrics
```bash
./starlink_exporter \
  -mode=web \
  -source=live \
  -listen=:9817 \
  -address=192.168.100.1:9200
```

Visit `http://localhost:9817/metrics` to view metrics.

#### Dummy Metrics (Testing)
```bash
./starlink_exporter \
  -mode=web \
  -source=dummy \
  -listen=:9817
```

Useful for testing without a physical Starlink dish.

### 3. Pushgateway Mode (Periodic Push)

#### Live Metrics
```bash
./starlink_exporter \
  -mode=pushgateway \
  -source=live \
  -pushgateway=http://pushgateway:9091 \
  -job=starlink_exporter \
  -instance=rpi-starlink-01 \
  -interval=15s \
  -address=192.168.100.1:9200
```

#### Dummy Metrics
```bash
./starlink_exporter \
  -mode=pushgateway \
  -source=dummy \
  -pushgateway=http://pushgateway:9091 \
  -job=starlink_exporter \
  -instance=rpi-starlink-01 \
  -interval=15s
```

## Command Line Options

```
-mode string
    Delivery mode: "web" or "pushgateway" (default "web")

-source string
    Metrics source: "live" (real dish) or "dummy" (simulated) (default "live")

-listen string
    Listen address for web mode (default ":9817")
    Example: :9817, 0.0.0.0:9817, 127.0.0.1:9817

-address string
    IP:port of Starlink dish (live mode only)
    Default: 192.168.100.1:9200

-pushgateway string
    Pushgateway URL (pushgateway mode only)
    Example: http://192.168.1.10:9091

-job string
    Job name for pushgateway grouping (default "starlink_exporter")

-instance string
    Instance name for pushgateway grouping (default "starlink_dish")

-interval string
    Push interval for pushgateway mode (default "15s")
    Example: 5s, 30s, 1m, 5m

-log-level string
    Log level: debug, info, warn, error (default "info")
```

## Docker Usage

### Docker Compose (Full Stack)

```bash
# Start full monitoring stack (Prometheus + Pushgateway + Grafana + Exporter)
docker-compose up -d

# Access:
# Grafana: http://localhost:3000 (admin/admin)
# Prometheus: http://localhost:9090
# Pushgateway: http://localhost:9091
# Exporter (web mode): http://localhost:9817/metrics
```

### Docker Run - Web Mode

```bash
docker run -d \
  --name starlink_exporter \
  -p 9817:9817 \
  -e SOURCE=dummy \
  ghcr.io/danopstech/starlink_exporter:latest
```

### Docker Run - Pushgateway Mode

```bash
docker run -d \
  --name starlink_exporter_push \
  -e MODE=pushgateway \
  -e SOURCE=dummy \
  -e PUSHGATEWAY=http://pushgateway:9091 \
  -e INTERVAL=15s \
  ghcr.io/danopstech/starlink_exporter:latest
```

### Docker Environment Variables

```
MODE=web|pushgateway
SOURCE=live|dummy
LISTEN=:9817
ADDRESS=192.168.100.1:9200
PUSHGATEWAY=http://pushgateway:9091
JOB=starlink_exporter
INSTANCE=starlink_docker
INTERVAL=15s
LOG_LEVEL=info
```

## Systemd Installation

### Web Mode

```bash
# Copy binary
sudo cp bin/starlink_exporter /usr/local/bin/

# Copy service file
sudo cp systemd/starlink_exporter_web.service /etc/systemd/system/

# Create user
sudo useradd -r -s /bin/false starlink || true

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable starlink_exporter_web.service
sudo systemctl start starlink_exporter_web.service

# Check status
sudo systemctl status starlink_exporter_web.service
```

### Pushgateway Mode

```bash
sudo cp systemd/starlink_exporter_pushgateway.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable starlink_exporter_pushgateway.service
sudo systemctl start starlink_exporter_pushgateway.service
```

## Prometheus Configuration

### Scrape Web Mode

```yaml
scrape_configs:
  - job_name: 'starlink_web'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:9817']
    metrics_path: '/metrics'
```

### Scrape Pushgateway

```yaml
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['localhost:9091']
```

## Dummy Metrics

When running with `-source=dummy`, the exporter generates realistic metrics:

- **Download**: 50-250 Mbps
- **Upload**: 5-40 Mbps
- **Latency**: 20-120 ms
- **Packet Loss**: 0-5%
- **SNR**: 5-15 dB
- **Obstruction**: 0-20%

Metrics update every 5 seconds with realistic variations.

## Exported Metrics

The exporter provides comprehensive Starlink metrics:

### Connection Health
- `starlink_dish_up` - Dish connection status
- `starlink_dish_pop_ping_latency_seconds` - Latency to PoP
- `starlink_dish_pop_ping_drop_ratio` - Packet loss ratio
- `starlink_dish_snr` - Signal-to-noise ratio

### Throughput
- `starlink_dish_downlink_throughput_bytes` - Download speed (bytes/sec)
- `starlink_dish_uplink_throughput_bytes` - Upload speed (bytes/sec)

### Obstructions
- `starlink_dish_currently_obstructed` - Currently obstructed (0/1)
- `starlink_dish_fraction_obstruction_ratio` - Obstruction percentage
- `starlink_dish_last_24h_obstructed_seconds` - Obstruction time in 24h
- `starlink_dish_wedge_fraction_obstruction_ratio` - Per-wedge obstruction

### Device Info
- `starlink_dish_info` - Device metadata (ID, version, country)
- `starlink_dish_uptime_seconds` - Device uptime

### Alerts
- `starlink_dish_alert_motors_stuck` - Motor stuck alert
- `starlink_dish_alert_thermal_throttle` - Thermal throttle alert
- `starlink_dish_alert_thermal_shutdown` - Thermal shutdown alert
- And more...

See full metric list in [metrics.md](./METRICS.md)

## Grafana Dashboard

Pre-built Grafana dashboard is included and automatically provisioned in docker-compose.

Dashboard features:
- Real-time throughput graphs
- Latency monitoring
- Signal strength tracking
- Obstruction visualization
- Packet loss monitoring
- Connection status indicator

## Development

### Project Structure

```
.
├── cmd/
│   └── starlink_exporter/
│       └── main.go
├── internal/
│   ├── provider/
│   │   ├── provider.go       (interface)
│   │   ├── live.go           (live provider)
│   │   └── dummy.go          (dummy provider)
│   ├── delivery/
│   │   ├── delivery.go       (interface)
│   │   ├── web.go            (web delivery)
│   │   └── pushgateway.go    (pushgateway delivery)
│   └── exporter/
│       └── exporter.go       (original exporter)
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── systemd/
│   ├── starlink_exporter_web.service
│   └── starlink_exporter_pushgateway.service
└── docker/
    ├── prometheus.yml
    └── grafana/
```

### Running Tests

```bash
make test

# With coverage
make test
open coverage.html
```

### Code Quality

```bash
make fmt
make lint
```

## Cross-Compilation

Build for multiple platforms:

```bash
# All platforms
make build-all

# Specific platforms
make build-amd64    # Linux x86_64
make build-arm64    # Linux ARM64 (RPi 4)
make build-armv7    # Linux ARMv7 (RPi 3, Zero)
```

Binaries will be in `dist/` directory.

## Performance

- Dummy provider generates metrics with negligible CPU usage
- Live provider minimal gRPC overhead
- Web mode: typical ~1-5 ms response time
- Pushgateway mode: sub-second push operation

## Troubleshooting

### Cannot connect to dish (live mode)

```bash
# Check dish is reachable
ping 192.168.100.1

# Test gRPC connection
grpcurl -plaintext 192.168.100.1:9200 describe
```

### Metrics not appearing in Prometheus

1. Check exporter is running: `curl http://localhost:9817/metrics`
2. Check Prometheus scrape config is correct
3. Check Prometheus logs: `curl http://localhost:9090/api/v1/status/config`

### Memory usage increasing

- Dummy provider includes per-5s updates, memory is stable
- Check for metric cardinality explosion (high-cardinality labels)

### Cannot push to Pushgateway

```bash
# Test pushgateway connectivity
curl -X POST http://pushgateway:9091/metrics/job/test

# Check logs
docker logs starlink_exporter
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

See [LICENSE](/LICENSE) file for details.

## Related Projects

- [Starlink Complete System](https://github.com/danopstech/starlink) - Full monitoring stack
- [Speedtest Exporter](https://github.com/danopstech/speedtest_exporter) - Speed test metrics
- [Prometheus](https://prometheus.io/) - Metrics collection
- [Grafana](https://grafana.com/) - Visualization

## Support

For issues, questions, or suggestions:

1. Check [existing issues](https://github.com/danopstech/starlink_exporter/issues)
2. Create a [new issue](https://github.com/danopstech/starlink_exporter/issues/new)
3. Include logs: `docker logs starlink_exporter`
4. Include command line flags used

---

**Note**: This is an enhanced version of the original Starlink exporter with additional delivery modes and data sources while maintaining backward compatibility.
