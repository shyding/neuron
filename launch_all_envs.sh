#!/bin/bash
set -e

echo "=== Cleaning old neuronex processes ==="
killall -9 neuronex 2>/dev/null || true
sleep 1

echo "=== Starting Port 8085 (Official NeuronEX) ==="
cd /opt/neuronex-4.0.2
nohup ./bin/neuronex alone > neuronex.log 2>&1 &
echo "Official NeuronEX started."
sleep 2

echo "=== Starting Port 8086 (Custom NeuronEX Standalone) ==="
chmod +x /root/neuronex_8086/bin/neuronex
cd /root/neuronex_8086
nohup ./bin/neuronex alone -c etc/neuronex.yaml > neuronex.log 2>&1 &
echo "Custom NeuronEX 8086 started."
sleep 2

echo "=== Checking Processes ==="
ps aux | grep -E "neuron|neuronex|kuiper" | grep -v grep

echo "=== Testing Endpoints ==="
curl -s -o /dev/null -w "Port 7000 status: %{http_code}\n" http://127.0.0.1:7000/api/v2/node -H "Authorization: Bearer test" || true
curl -s -o /dev/null -w "Port 8085 status: %{http_code}\n" http://127.0.0.1:8085/web/common/ || true
curl -s -o /dev/null -w "Port 8086 status: %{http_code}\n" http://127.0.0.1:8086/web/common/ || true
