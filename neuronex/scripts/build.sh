#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

echo "=== [1/2] Building NeuronEX Gateway (Go) ==="
cd "${ROOT_DIR}/neuronex"
go build -ldflags="-s -w" -o "${ROOT_DIR}/build/bin/neuronex" ./cmd/neuronex

echo "=== [2/2] Preparing NeuronEX Distribution Tree ==="
mkdir -p "${ROOT_DIR}/build/web/common"
mkdir -p "${ROOT_DIR}/build/etc"
mkdir -p "${ROOT_DIR}/build/run/neuron"
mkdir -p "${ROOT_DIR}/build/log/neuronex"

# Copy Web Assets and Configuration
rsync -a "${ROOT_DIR}/neuronex/web/common/" "${ROOT_DIR}/build/web/common/"
cp "${ROOT_DIR}/neuronex/etc/neuronex.yaml" "${ROOT_DIR}/build/etc/neuronex.yaml"

echo "=== NeuronEX Build Complete: ${ROOT_DIR}/build/bin/neuronex ==="
