# Installation Guide

## Prerequisites

- Go 1.21+ (for building from source)
- Docker & Docker Compose (for containerized deployment)
- Starlink dish running Starlink OS (for live mode)
- Linux system for systemd deployment

## Table of Contents

- [Binary Installation](#binary-installation)
- [Docker Installation](#docker-installation)
- [Raspberry Pi Setup](#raspberry-pi-setup)
- [Systemd Service](#systemd-service-installation)
- [Configuration Examples](#configuration-examples)
- [Network Setup](#network-setup)

## Binary Installation

### Download Pre-built Binary

Download the latest binary for your platform from [releases](https://github.com/danopstech/starlink_exporter/releases):

```bash
# For Linux x86_64
wget https://github.com/danopstech/starlink_exporter/releases/download/v1.0.0/starlink_exporter-linux-amd64
chmod +x starlink_exporter-linux-amd64
./starlink_exporter-linux-amd64 -help

# For Raspberry Pi 4 (ARM64)
wget https://github.com/danopstech/starlink_exporter/releases/download/v1.0.0/starlink_exporter-linux-arm64
chmod +x starlink_exporter-linux-arm64
./starlink_exporter-linux-arm64 -help

# For Raspberry Pi 3/Zero (ARMv7)
wget https://github.com/danopstech/starlink_exporter/releases/download/v1.0.0/starlink_exporter-linux-armv7
chmod +x starlink_exporter-linux-armv7
./starlink_exporter-linux-armv7 -help
```

### Build from Source

```bash
git clone https://github.com/danopstech/starlink_exporter.git
cd starlink_exporter

# Build for current platform
make build

# Binary is in ./bin/starlink_exporter
./bin/starlink_exporter -help
```

## Docker Installation

### Quick Start with Docker

```bash
# Web mode with dummy metrics
docker run -d \
  --name starlink_exporter \
  -p 9817:9817 \
  -e SOURCE=dummy \
  -e MODE=web \
  ghcr.io/danopstech/starlink_exporter:latest

# View metrics
curl http://localhost:9817/metrics
```

### Build Docker Image Locally

```bash
git clone https://github.com/danopstech/starlink_exporter.git
cd starlink_exporter

make build-docker

# Run locally built image
docker run -d \
  --name starlink_exporter \
  -p 9817:9817 \
  starlink_exporter:latest
```

## Docker Compose Installation

### Complete Monitoring Stack

```bash
git clone https://github.com/danopstech/starlink_exporter.git
cd starlink_exporter

# Start all services (Prometheus, Pushgateway, Grafana, Exporter)
docker-compose up -d

# Access services:
# Grafana: http://localhost:3000 (admin/admin)
# Prometheus: http://localhost:9090
# Pushgateway: http://localhost:9091
# Exporter: http://localhost:9817/metrics

# View logs
docker-compose logs -f exporter_web
docker-compose logs -f exporter_pushgateway
```

### Customize Docker Compose

Edit `docker-compose.yml` to change:

```yaml
environment:
  - MODE=web|pushgateway
  - SOURCE=live|dummy
  - LISTEN=:9817
  - PUSHGATEWAY=http://pushgateway:9091
  - INTERVAL=15s
```

## Raspberry Pi Setup

### Installation on Raspberry Pi 4 (ARM64)

```bash
# 1. Download ARM64 binary
wget https://github.com/danopstech/starlink_exporter/releases/download/v1.0.0/starlink_exporter-linux-arm64

# 2. Make executable
chmod +x starlink_exporter-linux-arm64

# 3. Test run
./starlink_exporter-linux-arm64 -mode=web -source=dummy -listen=:9817

# 4. Copy to system location
sudo cp starlink_exporter-linux-arm64 /usr/local/bin/starlink_exporter
```

### Installation on Raspberry Pi 3/Zero (ARMv7)

```bash
# Download ARMv7 binary
wget https://github.com/danopstech/starlink_exporter/releases/download/v1.0.0/starlink_exporter-linux-armv7
chmod +x starlink_exporter-linux-armv7
sudo cp starlink_exporter-linux-armv7 /usr/local/bin/starlink_exporter
```

### Raspberry Pi Network Setup

```bash
# Find Starlink dish IP (usually on same network)
ping 192.168.100.1

# Test gRPC connection
sudo apt-get install grpcurl
grpcurl -plaintext 192.168.100.1:9200 describe

# Configure firewall if needed
sudo ufw allow 9817/tcp
sudo ufw allow 9091/tcp
```

## Systemd Service Installation

### Web Mode Service

```bash
# 1. Create user
sudo useradd -r -s /bin/false starlink

# 2. Copy binary
sudo cp bin/starlink_exporter /usr/local/bin/

# 3. Copy service file
sudo cp systemd/starlink_exporter_web.service /etc/systemd/system/

# 4. Reload systemd
sudo systemctl daemon-reload

# 5. Enable service (auto-start)
sudo systemctl enable starlink_exporter_web.service

# 6. Start service
sudo systemctl start starlink_exporter_web.service

# 7. Check status
sudo systemctl status starlink_exporter_web.service

# 8. View logs
sudo journalctl -u starlink_exporter_web.service -f
```

### Pushgateway Mode Service

```bash
# Same as above but use:
sudo cp systemd/starlink_exporter_pushgateway.service /etc/systemd/system/
sudo systemctl enable starlink_exporter_pushgateway.service
sudo systemctl start starlink_exporter_pushgateway.service
```

### Custom Systemd Service

Create `/etc/systemd/system/starlink_exporter_custom.service`:

```ini
[Unit]
Description=Starlink Exporter - Custom Configuration
After=network.target

[Service]
Type=simple
User=starlink
Group=starlink
ExecStart=/usr/local/bin/starlink_exporter \
  -mode=pushgateway \
  -source=live \
  -address=192.168.100.1:9200 \
  -pushgateway=http://prometheus-server:9091 \
  -job=starlink \
  -instance=my-starlink \
  -interval=30s \
  -log-level=debug

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Then:

```bash
sudo systemctl daemon-reload
sudo systemctl enable starlink_exporter_custom.service
sudo systemctl start starlink_exporter_custom.service
```

## Configuration Examples

### Example 1: Web Mode on Raspberry Pi with Live Data

```bash
./starlink_exporter \
  -mode=web \
  -source=live \
  -listen=0.0.0.0:9817 \
  -address=192.168.100.1:9200 \
  -log-level=info
```

**Prometheus config:**

```yaml
scrape_configs:
  - job_name: 'starlink'
    static_configs:
      - targets: ['raspberrypi:9817']
```

### Example 2: Pushgateway with Prometheus

```bash
./starlink_exporter \
  -mode=pushgateway \
  -source=live \
  -pushgateway=http://prometheus-server:9091 \
  -job=starlink_exporter \
  -instance=rpi-starlink-01 \
  -interval=15s \
  -address=192.168.100.1:9200
```

**Prometheus config:**

```yaml
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['prometheus-server:9091']
```

### Example 3: Multiple Instances with Pushgateway

```bash
# Instance 1
./starlink_exporter \
  -mode=pushgateway \
  -source=live \
  -pushgateway=http://pushgateway:9091 \
  -job=starlink \
  -instance=location-1 \
  -address=192.168.100.1:9200

# Instance 2 (different address)
./starlink_exporter \
  -mode=pushgateway \
  -source=live \
  -pushgateway=http://pushgateway:9091 \
  -job=starlink \
  -instance=location-2 \
  -address=192.168.50.1:9200
```

**Prometheus config:** Same job_name scrapes both instances from pushgateway

### Example 4: Testing with Dummy Metrics

```bash
# Start exporter with dummy metrics
./starlink_exporter \
  -mode=web \
  -source=dummy \
  -listen=:9817

# In another terminal, test with Prometheus
docker run -d \
  -p 9090:9090 \
  -v prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

## Network Setup

### Finding Your Starlink Dish

```bash
# Check if dish is on network
ping 192.168.100.1

# Find all devices on network
nmap -sn 192.168.100.0/24

# Test gRPC port
nc -zv 192.168.100.1 9200
```

### Firewall Configuration

#### UFW (Ubuntu/Debian)

```bash
# Allow exporter web port
sudo ufw allow 9817/tcp

# Allow from specific IP
sudo ufw allow from 192.168.1.100 to any port 9817

# Allow pushgateway port
sudo ufw allow 9091/tcp
```

#### iptables

```bash
# Allow TCP 9817
sudo iptables -A INPUT -p tcp --dport 9817 -j ACCEPT

# Allow from specific network
sudo iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 9817 -j ACCEPT

# Save rules
sudo iptables-save > /etc/iptables/rules.v4
```

### Port Forwarding (Optional)

If accessing from remote network:

```bash
# SSH tunnel to Raspberry Pi
ssh -L 9817:localhost:9817 user@raspberry-pi-ip

# Then access locally
curl http://localhost:9817/metrics
```

## Verification

### Check Installation

```bash
# Check version
./starlink_exporter -help

# Test connection
./starlink_exporter -mode=web -source=dummy &
sleep 2
curl http://localhost:9817/metrics | head -20

# Kill test process
pkill -f starlink_exporter
```

### Verify Metrics Collection

```bash
# Web mode
curl http://localhost:9817/metrics | grep starlink_dish_up

# Should return:
# starlink_dish_up 1
```

### Monitor Logs

```bash
# Systemd service
sudo journalctl -u starlink_exporter_web.service -f

# Docker
docker logs -f starlink_exporter

# Docker Compose
docker-compose logs -f exporter_web
```

## Upgrading

### Binary

```bash
# Download new version
wget https://github.com/danopstech/starlink_exporter/releases/download/vX.Y.Z/starlink_exporter-linux-amd64

# Stop running instance
sudo systemctl stop starlink_exporter_web.service

# Backup old binary
sudo cp /usr/local/bin/starlink_exporter /usr/local/bin/starlink_exporter.bak

# Copy new binary
sudo cp starlink_exporter-linux-amd64 /usr/local/bin/starlink_exporter

# Start new version
sudo systemctl start starlink_exporter_web.service
```

### Docker

```bash
docker-compose pull
docker-compose up -d --force-recreate
```

## Troubleshooting Installation

### Binary not found
```bash
# Ensure binary is in PATH
echo $PATH
sudo cp starlink_exporter /usr/local/bin/
```

### Permission denied
```bash
chmod +x ./starlink_exporter
sudo chown starlink:starlink /usr/local/bin/starlink_exporter
```

### Cannot connect to Starlink dish
```bash
# Test network connectivity
ping 192.168.100.1

# Test gRPC
grpcurl -plaintext 192.168.100.1:9200 describe

# Check firewall
sudo iptables -L | grep 9200
```

### Systemd service fails
```bash
# Check service status
sudo systemctl status starlink_exporter_web.service

# View detailed logs
sudo journalctl -u starlink_exporter_web.service -n 50

# Verify paths in service file
sudo cat /etc/systemd/system/starlink_exporter_web.service
```
