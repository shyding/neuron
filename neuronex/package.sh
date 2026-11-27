#!/bin/bash
DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"

if command -v python3 &>/dev/null; then
    exec python3 package.py "$@"
elif command -v python &>/dev/null; then
    exec python package.py "$@"
else
    echo "Python 3 is required to run the packaging script."
    exit 1
fi
