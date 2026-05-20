#!/bin/sh
set -eu

VERSION="v1.0.0"
BINARY_URL="https://github.com/adityazxzxz/starlink-exporter/releases/download/${VERSION}/starlink-exporter-arm64"
INSTALL_PATH="/usr/local/bin/starlink-exporter"
SERVICE_FILE="/etc/systemd/system/starlink-exporter@.service"
DEFAULT_ENV="/etc/default/starlink-exporter"

if [ "$(id -u)" -ne 0 ]; then
    echo "Please run as root (sudo)."
    exit 1
fi

command -v curl >/dev/null 2>&1 || {
    echo "curl is required"
    exit 1
}

command -v systemctl >/dev/null 2>&1 || {
    echo "systemd is required"
    exit 1
}

echo "=== Starlink Exporter Installer ==="
echo

printf "Instance name (example: rumah, gudang-jakarta, dish-01): "
read INSTANCE

if [ -z "$INSTANCE" ]; then
    echo "Instance name cannot be empty"
    exit 1
fi

printf "Pushgateway URL [http://pushgateway.example.com:9091]: "
read PUSHGATEWAY
PUSHGATEWAY=${PUSHGATEWAY:-http://pushgateway.example.com:9091}

printf "Job name [starlink]: "
read JOB
JOB=${JOB:-starlink}

printf "Interval [30s]: "
read INTERVAL
INTERVAL=${INTERVAL:-30s}

echo
echo "Downloading binary..."
curl -fsSL "$BINARY_URL" -o "$INSTALL_PATH"
chmod +x "$INSTALL_PATH"

echo "Creating default config..."
cat > "$DEFAULT_ENV" <<EOF
MODE=pushgateway
SOURCE=live
EOF

echo "Creating instance config..."
cat > "/etc/default/starlink-exporter-${INSTANCE}" <<EOF
PUSHGATEWAY=${PUSHGATEWAY}
JOB=${JOB}
INTERVAL=${INTERVAL}
EOF

echo "Creating systemd service..."
cat > "$SERVICE_FILE" <<'EOF'
[Unit]
Description=Starlink Exporter (%i)
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root

EnvironmentFile=-/etc/default/starlink-exporter
EnvironmentFile=-/etc/default/starlink-exporter-%i

ExecStart=/usr/local/bin/starlink-exporter \
  -mode=${MODE} \
  -source=${SOURCE} \
  -pushgateway=${PUSHGATEWAY} \
  -job=${JOB} \
  -instance=%i \
  -interval=${INTERVAL}

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd..."
systemctl daemon-reload

echo "Enabling service..."
systemctl enable --now "starlink-exporter@${INSTANCE}"

echo
echo "Installation complete."
echo
echo "Service name:"
echo "  starlink-exporter@${INSTANCE}"
echo
echo "Check status:"
echo "  systemctl status starlink-exporter@${INSTANCE}"
echo
echo "View logs:"
echo "  journalctl -u starlink-exporter@${INSTANCE} -f"