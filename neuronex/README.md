# NeuronEX (Self-Contained Implementation)

NeuronEX is a unified IIoT Edge Platform combining:
1. **Neuron IIoT Core**: High-concurrency industrial protocol data acquisition (Modbus, Siemens, Mitsubishi, Omron, Rockwell, CNC, IEC 104/61850, etc.).
2. **eKuiper**: LF Edge lightweight streaming data processing engine and SQL rule engine.
3. **NeuronEX Gateway (Go)**: Reverse proxy, supervisor process manager, JWT authentication (admin/0000), and static web asset delivery.
4. **Unified Web Dashboard**: Comprehensive Vue 3 / Element Plus SPA management interface.

---

## 1. Directory Structure

```text
neuronex/
├── cmd/
│   └── neuronex/
│       └── main.go           # CLI entrypoint (run, start, alone, reset-password)
├── internal/
│   ├── auth/                 # JWT token generation, verification, and credentials
│   ├── config/               # neuronex.yaml configuration loader
│   ├── proxy/                # High-performance reverse proxy for Neuron & eKuiper
│   ├── server/               # HTTP server & static web assets router
│   └── supervisor/           # Child process launcher and supervisor (Neuron, eKuiper)
├── etc/
│   └── neuronex.yaml         # Default configuration (ports, reverse proxy mappings)
├── web/
│   └── common/               # Unified Web Dashboard SPA frontend
├── scripts/
│   ├── build.sh              # Build script for NeuronEX
│   └── run.sh                # Launcher script for NeuronEX
└── go.mod
```

---

## 2. Build Instructions

Ensure Go 1.20+ is installed:

```bash
cd neuronex
go build -ldflags="-s -w" -o bin/neuronex ./cmd/neuronex
```

Or run the automated build script:

```bash
bash neuronex/scripts/build.sh
```

---

## 3. Running NeuronEX

Run in console mode:

```bash
cd build  # or project root
./bin/neuronex run -c etc/neuronex.yaml
```

- **Web Dashboard URL**: `http://localhost:8085/web/common/`
- **Default Username**: `admin`
- **Default Password**: `0000`
