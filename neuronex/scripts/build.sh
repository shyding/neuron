#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

echo "=== [1/3] Building NeuronEX Gateway (Go) ==="
cd "${ROOT_DIR}/neuronex"
go build -ldflags="-s -w" -o "${ROOT_DIR}/build/bin/neuronex" ./cmd/neuronex

echo "=== [2/3] Building Frontend Dashboard (Vue 3 / TypeScript) ==="
if command -v yarn >/dev/null 2>&1 && [ -f "${ROOT_DIR}/dashboard/package.json" ]; then
    echo "Building dashboard with yarn..."
    cd "${ROOT_DIR}/dashboard"
    yarn build || true
fi

echo "=== [3/3] Preparing NeuronEX Distribution Tree ==="
mkdir -p "${ROOT_DIR}/build/web/common"
mkdir -p "${ROOT_DIR}/build/etc"
mkdir -p "${ROOT_DIR}/build/run/neuron"
mkdir -p "${ROOT_DIR}/build/log/neuronex"

# Copy Web Assets and Configuration
if [ -d "${ROOT_DIR}/dashboard/dist" ]; then
    rsync -a "${ROOT_DIR}/dashboard/dist/" "${ROOT_DIR}/build/web/common/"
else
    rsync -a "${ROOT_DIR}/neuronex/web/common/" "${ROOT_DIR}/build/web/common/"
fi
cp "${ROOT_DIR}/neuronex/etc/neuronex.yaml" "${ROOT_DIR}/build/etc/neuronex.yaml"

echo "=== NeuronEX Full-Stack Build Complete: ${ROOT_DIR}/build/bin/neuronex ==="
