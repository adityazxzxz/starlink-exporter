# Starlink Exporter

Exporter Prometheus standalone untuk mengumpulkan metrik dari Starlink Dish melalui gRPC. Mendukung dua mode pengiriman: **Web** (untuk scraping langsung oleh Prometheus) dan **Pushgateway** (untuk push metrik secara berkala).

## Build

Pastikan Anda memiliki Go 1.16 atau lebih tinggi.

```bash
# Build untuk platform saat ini
go build -o starlink-exporter

# Build untuk ARM64 (Raspberry Pi 4)
GOOS=linux GOARCH=arm64 go build -o starlink-exporter-arm64

# Build untuk ARMv7 (Raspberry Pi 3/Zero)
GOOS=linux GOARCH=arm go build -o starlink-exporter-armv7
```

## Penggunaan

### Mode Web (Scraping Langsung)

Jalankan exporter dengan mode web (default) dan sumber live:

```bash
./starlink-exporter -mode=web -source=live
```

Prometheus akan scrape metrik di `http://localhost:9817/metrics`

**Contoh dengan custom listen address:**

```bash
./starlink-exporter -mode=web -source=live -listen=:8080
```

**Menggunakan dummy data (tanpa Starlink dish):**

```bash
./starlink-exporter -mode=web -source=dummy
```

### Mode Pushgateway

Push metrik secara berkala ke Pushgateway:

```bash
./starlink-exporter -mode=pushgateway -source=live \
  -pushgateway=http://localhost:9091
```

**Dengan custom job dan instance:**

```bash
./starlink-exporter -mode=pushgateway -source=live \
  -pushgateway=http://pushgateway.example.com:9091 \
  -job=starlink \
  -instance=my-dish \
  -interval=30s
```

**Menggunakan dummy data dengan pushgateway:**

```bash
./starlink-exporter -mode=pushgateway -source=dummy \
  -pushgateway=http://localhost:9091
```

## Command Line Arguments

| Flag | Default | Deskripsi |
|------|---------|-----------|
| `-mode` | `web` | Delivery mode: `web` (HTTP server) atau `pushgateway` (push ke Pushgateway) |
| `-source` | `live` | Metrics source: `live` (dari Starlink dish) atau `dummy` (data simulasi) |
| `-listen` | `:9817` | Listen address untuk web mode (contoh: `:9817`, `0.0.0.0:8080`) |
| `-address` | `192.168.100.1:50051` | IP address dan port untuk reach Starlink dish (live mode only) |
| `-pushgateway` | `` | Pushgateway URL (required untuk pushgateway mode, contoh: `http://localhost:9091`) |
| `-job` | `starlink_exporter` | Job name untuk pushgateway |
| `-instance` | `starlink_dish` | Instance name untuk pushgateway |
| `-interval` | `15s` | Push interval untuk pushgateway mode (contoh: `30s`, `1m`) |
| `-log-level` | `info` | Log level: `debug`, `info`, `warn`, `error` |

## Contoh Penggunaan

### 1. Web Mode dengan Dish Asli

```bash
./starlink-exporter -mode=web -source=live -address=192.168.100.1:50051
```

Akses metrik di: `http://localhost:9817/metrics`

### 2. Web Mode dengan Dummy Data

```bash
./starlink-exporter -mode=web -source=dummy
```

Berguna untuk testing dan development tanpa hardware Starlink.

### 3. Pushgateway Mode dengan Dish Asli

```bash
./starlink-exporter -mode=pushgateway -source=live \
  -pushgateway=http://prometheus-stack.local:9091 \
  -job=starlink \
  -interval=60s
```

### 4. Pushgateway Mode dengan Dummy Data

```bash
./starlink-exporter -mode=pushgateway -source=dummy \
  -pushgateway=http://localhost:9091 \
  -interval=30s
```

### 5. Debug dengan Log Level Verbose

```bash
./starlink-exporter -mode=web -source=live -log-level=debug
```

## Metrik

Exporter mengumpulkan metrik Starlink berikut:

- **Download speed** (bytes/sec)
- **Upload speed** (bytes/sec)
- **Latency** (ms)
- **Packet loss** (%)
- **Signal-to-noise ratio** (dB)
- **Obstruction percentage** (%)
- **Uptime** (seconds)
- **Connection state** dan alerts

## Health Check

Untuk web mode, endpoint `/health` tersedia untuk memeriksa status koneksi gRPC ke Starlink dish:

```bash
curl http://localhost:9817/health
```



## Installation

```bash
wget -O install.sh https://github.com/adityazxzxz/starlink-exporter/releases/download/v1.0.1/install.sh

chmod +x install.sh
sudo ./install.sh
```
