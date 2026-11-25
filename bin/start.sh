#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

chmod +x ./neuron 2>/dev/null || true
chmod +x ./modbus_simulator 2>/dev/null || true

export LD_LIBRARY_PATH="$DIR:$DIR/plugins:$LD_LIBRARY_PATH"
echo "Starting Neuron IIoT Server from $DIR..."
echo "Web dashboard will be available at: http://localhost:7000"
exec ./neuron
