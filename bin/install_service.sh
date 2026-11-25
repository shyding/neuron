#!/bin/bash
set -e

if [ "$EUID" -ne 0 ]; then
  echo "Error: Please run as root (e.g. sudo ./install_service.sh)"
  exit 1
fi

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chmod +x "$DIR/neuron" "$DIR/start.sh" 2>/dev/null || true

SERVICE_FILE="/etc/systemd/system/neuron.service"

cat <<EOF > "$SERVICE_FILE"
[Unit]
Description=Neuron IIoT Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DIR
Environment="LD_LIBRARY_PATH=$DIR:$DIR/plugins"
ExecStart=$DIR/neuron
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable neuron
systemctl restart neuron

echo "Neuron service installed and started successfully!"
echo "Check status with: systemctl status neuron"
echo "Web dashboard: http://localhost:7000"
