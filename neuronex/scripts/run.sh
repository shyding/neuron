#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Check if running from build directory or installed directory
if [ -f "${ROOT_DIR}/build/bin/neuronex" ]; then
    DIST_DIR="${ROOT_DIR}/build"
elif [ -f "${ROOT_DIR}/bin/neuronex" ]; then
    DIST_DIR="${ROOT_DIR}"
else
    echo "NeuronEX binary not found. Running build first..."
    "${SCRIPT_DIR}/build.sh"
    DIST_DIR="${ROOT_DIR}/build"
fi

cd "${DIST_DIR}"
echo "Starting NeuronEX Platform from ${DIST_DIR}..."
echo "Web UI will be available at: http://localhost:8085/web/common/"
exec ./bin/neuronex run -c etc/neuronex.yaml
