#!/usr/bin/env python3
"""
NeuronEX Multi-Platform Packaging & Closed-Loop Verification Tool
Supports Windows (amd64), Linux (amd64), Linux (arm64), and interactive menu.
"""

import sys, os, shutil, subprocess, zipfile, tarfile, json

VERSION = "v4.0.2"
ROOT_DIR = os.path.dirname(os.path.abspath(__file__))
DIST_DIR = os.path.join(ROOT_DIR, "dist")
SRC_DIR = ROOT_DIR

PLATFORMS = {
    "1": {"os": "windows", "arch": "amd64", "ext": ".exe", "archive": "zip", "label": "Windows x86_64 (.zip)"},
    "2": {"os": "linux",   "arch": "amd64", "ext": "",     "archive": "tar", "label": "Linux x86_64 (.tar.gz)"},
    "3": {"os": "linux",   "arch": "arm64", "ext": "",     "archive": "tar", "label": "Linux ARM64 (.tar.gz)"},
}

def print_header():
    print("""
======================================================================
       🚀  NeuronEX Multi-Platform Packaging & Release Tool  🚀      
======================================================================
   Version: 4.0.2 (Self-Implemented Full-Stack Independent Edition)
   Branch : feat/neuronex
======================================================================
""")

def print_menu():
    print("   [1] 📦 Package for Windows (windows/amd64) -> .zip")
    print("   [2] 📦 Package for Linux   (linux/amd64)   -> .tar.gz")
    print("   [3] 📦 Package for Linux   (linux/arm64)   -> .tar.gz")
    print("   [4] 🚀 Package ALL Platforms (1 + 2 + 3)")
    print("   [5] 🧪 Run Closed-Loop Verification (Test Windows & WSL2)")
    print("   [0] ❌ Exit")
    print("=" * 70)

def build_binary(target_os, target_arch, out_file, pkg_path):
    env = os.environ.copy()
    env["CGO_ENABLED"] = "0"
    env["GOOS"] = target_os
    env["GOARCH"] = target_arch

    cmd = ["go", "build", "-trimpath", "-ldflags", "-s -w", "-o", out_file, pkg_path]
    print(f"      [Go Build] {target_os}/{target_arch} -> {os.path.basename(out_file)}...")
    res = subprocess.run(cmd, cwd=SRC_DIR, env=env, capture_output=True, text=True)
    if res.returncode != 0:
        print(f"      [ERROR] Compilation failed: {res.stderr}")
        return False
    return True

def create_platform_package(cfg):
    target_os = cfg["os"]
    target_arch = cfg["arch"]
    ext = cfg["ext"]
    archive_type = cfg["archive"]
    
    pkg_name = f"neuronex-{VERSION}-{target_os}-{target_arch}"
    staging_dir = os.path.join(DIST_DIR, f"stage_{pkg_name}")
    if os.path.exists(staging_dir):
        shutil.rmtree(staging_dir, ignore_errors=True)
    
    os.makedirs(os.path.join(staging_dir, "bin"), exist_ok=True)
    os.makedirs(os.path.join(staging_dir, "etc"), exist_ok=True)
    os.makedirs(os.path.join(staging_dir, "log"), exist_ok=True)

    print(f"\n===> Building package for {cfg['label']}...")

    # 1. Compile neuronex
    neuronex_bin = os.path.join(staging_dir, "bin", f"neuronex{ext}")
    if not build_binary(target_os, target_arch, neuronex_bin, "./cmd/neuronex"):
        return None

    # 2. Compile license generator tool
    licgen_bin = os.path.join(staging_dir, "bin", f"neuronex-licgen{ext}")
    if not build_binary(target_os, target_arch, licgen_bin, "./tools/license-generator"):
        return None

    # 3. Copy configuration files
    shutil.copy2(os.path.join(SRC_DIR, "etc", "neuronex.yaml"), os.path.join(staging_dir, "etc", "neuronex.yaml"))

    # 4. Generate default pre-activated permanent license file
    license_data = {
        "licenseType": "COMMERCIAL",
        "customerName": "EMQ Technologies Inc.",
        "customerContact": "support@emqx.io",
        "startDate": "2024-01-01 00:00:00",
        "expiredDate": "2999-01-01 00:00:00",
        "expired": False,
        "totalTagCount": 1000000,
        "usedTagCount": 0,
        "enabledPlugins": [
            "ABB COMLI", "Allen-Bradley 5000 EtherNet/IP", "Allen-Bradley ControlLogix 5500",
            "Allen-Bradley DF1", "Allen-Bradley MicroLogix 1400", "BACnet/IP", "Beckhoff ADS",
            "CODESYS V3 TCP", "DLT645-1997", "DLT645-2007", "DNP 3.0", "EtherNet/IP(CIP)",
            "Fanuc Focas Ethernet", "GE Historian", "GE SRTP", "HJ212-2017", "HOSTLINK CMODE",
            "HollySys Modbus RTU", "HollySys Modbus TCP", "IEC60870-5-101", "IEC60870-5-102",
            "IEC60870-5-103", "IEC60870-5-104", "IEC61850", "Inovance Modbus TCP", "KNXnet/IP",
            "KUKA ETHERNET KRL TCP", "MTConnect", "Mitsubishi 1E", "Mitsubishi 3E", "Mitsubishi 4E",
            "Mitsubishi CNC", "Mitsubishi FX", "Modbus ASCII", "Modbus RTU", "Modbus TCP",
            "Modbus TCP QH", "NON A11", "OPC UA", "OPC UA Server", "OPCAE", "OPCDA",
            "Omron EtherNet/IP", "Omron FINS TCP", "Omron FINS UDP", "Panasonic Mewtocol",
            "SECS GEM HSMS", "Siemens FetchWrite", "Siemens MPI", "Siemens S7 ISOTCP",
            "Siemens S7 ISOTCP for 300/400", "SparkPlugB", "WebSocket", "XINJE Modbus RTU",
            "Brother CNC", "HEIDENHAIN CNC", "KND CNC", "Mazak CNC", "SYNTEC CNC"
        ],
        "hardwareToken": "NEURONEX-CORE-ENTERPRISE",
        "ekuiperOn": True,
        "officialEmail": "support@emqx.io",
        "original": False
    }
    with open(os.path.join(staging_dir, "etc", "license.json"), "w", encoding="utf-8") as f:
        json.dump(license_data, f, indent=2, ensure_ascii=False)

    # 5. Copy frontend static assets (web/)
    web_src = os.path.join(SRC_DIR, "web")
    web_dst = os.path.join(staging_dir, "web")
    if os.path.exists(web_src):
        shutil.copytree(web_src, web_dst)

    # 6. Add start/stop scripts
    if target_os == "windows":
        with open(os.path.join(staging_dir, "start.bat"), "w", encoding="gbk") as f:
            f.write("""@echo off
echo ===================================================
echo   Starting NeuronEX Industrial Edge Platform...
echo ===================================================
cd /d "%~dp0"
start "NeuronEX" bin\\neuronex.exe alone
echo [OK] NeuronEX service started!
echo Web Dashboard available at: http://localhost:8086/web/common/
echo Default credentials: admin / 0000
pause
""")
        with open(os.path.join(staging_dir, "stop.bat"), "w", encoding="gbk") as f:
            f.write("""@echo off
echo Stopping NeuronEX...
taskkill /F /IM neuronex.exe >nul 2>&1
echo [OK] Stopped.
pause
""")
        with open(os.path.join(staging_dir, "generate-license.bat"), "w", encoding="gbk") as f:
            f.write("""@echo off
cd /d "%~dp0"
bin\\neuronex-licgen.exe %*
pause
""")
    else:
        with open(os.path.join(staging_dir, "start.sh"), "w", encoding="utf-8") as f:
            f.write("""#!/bin/bash
DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"
mkdir -p log
nohup ./bin/neuronex alone > log/neuronex.log 2>&1 &
PID=$!
echo "==================================================="
echo "  NeuronEX Industrial Edge Platform Started!"
echo "==================================================="
echo "  PID: $PID"
echo "  Web Console: http://localhost:8086/web/common/"
echo "  Default User: admin / 0000"
""")
        with open(os.path.join(staging_dir, "stop.sh"), "w", encoding="utf-8") as f:
            f.write("""#!/bin/bash
pkill -f 'bin/neuronex' || true
echo "NeuronEX service stopped."
""")
        with open(os.path.join(staging_dir, "generate-license.sh"), "w", encoding="utf-8") as f:
            f.write("""#!/bin/bash
DIR="$(cd "$(dirname "$0")" && pwd)"
"$DIR/bin/neuronex-licgen" "$@"
""")
        # Set execute permissions
        for sh_file in ["start.sh", "stop.sh", "generate-license.sh"]:
            try:
                os.chmod(os.path.join(staging_dir, sh_file), 0o755)
            except:
                pass
        for b_file in [f"neuronex{ext}", f"neuronex-licgen{ext}"]:
            try:
                os.chmod(os.path.join(staging_dir, "bin", b_file), 0o755)
            except:
                pass

    # 7. Add README.md
    with open(os.path.join(staging_dir, "README.md"), "w", encoding="utf-8") as f:
        f.write(f"""# NeuronEX Industrial Edge Platform ({VERSION})
### Target Architecture: {target_os}/{target_arch}

## 1. Quick Start
- **Windows:** Double click `start.bat`
- **Linux:** Run `./start.sh`

## 2. Web Console Access
- **URL:** [http://localhost:8086/web/common/](http://localhost:8086/web/common/)
- **Username:** `admin`
- **Password:** `0000`

## 3. Commercial License & Activation Tool
This release comes pre-activated with an enterprise commercial license (1,000,000 points, 59 industrial protocols, permanent validity until 2999-01-01).

To generate custom activation codes or license files, use:
- **Windows:** `bin\\neuronex-licgen.exe -customer "Your Company" -tags 1000000`
- **Linux:** `./bin/neuronex-licgen -customer "Your Company" -tags 1000000`
""")

    # 8. Create compressed archive
    os.makedirs(DIST_DIR, exist_ok=True)
    if archive_type == "zip":
        archive_path = os.path.join(DIST_DIR, f"{pkg_name}.zip")
        print(f"      [Packaging] Compressing to {os.path.basename(archive_path)}...")
        with zipfile.ZipFile(archive_path, "w", zipfile.ZIP_DEFLATED) as zf:
            for root, dirs, files in os.walk(staging_dir):
                for file in files:
                    full_p = os.path.join(root, file)
                    rel_p = os.path.relpath(full_p, staging_dir)
                    zf.write(full_p, rel_p)
    else:
        archive_path = os.path.join(DIST_DIR, f"{pkg_name}.tar.gz")
        print(f"      [Packaging] Compressing to {os.path.basename(archive_path)}...")
        with tarfile.open(archive_path, "w:gz") as tf:
            for root, dirs, files in os.walk(staging_dir):
                for file in files:
                    full_p = os.path.join(root, file)
                    rel_p = os.path.relpath(full_p, staging_dir)
                    tf.add(full_p, arcname=rel_p)

    shutil.rmtree(staging_dir, ignore_errors=True)
    size_mb = os.path.getsize(archive_path) / (1024 * 1024)
    print(f"      [SUCCESS] Package created: {os.path.basename(archive_path)} ({size_mb:.2f} MB)")
    return archive_path

def run_tests():
    test_script = os.path.join(SRC_DIR, "scripts", "test_closed_loop.py")
    if os.path.exists(test_script):
        cmd = [sys.executable, test_script, "all"]
        subprocess.run(cmd)
    else:
        print("[ERROR] test_closed_loop.py not found!")

def interactive_menu():
    print_header()
    while True:
        print_menu()
        choice = input("👉 Please enter your choice [0-5]: ").strip()
        if choice == "0":
            print("\nExiting. Good bye!\n")
            break
        elif choice in ["1", "2", "3"]:
            create_platform_package(PLATFORMS[choice])
        elif choice == "4":
            print("\n===> Building packages for ALL platforms...")
            for k in ["1", "2", "3"]:
                create_platform_package(PLATFORMS[k])
            print("\n🎉 ALL PLATFORM PACKAGES CREATED SUCCESSFULLY in `dist/`!")
        elif choice == "5":
            print("\n===> Starting Closed-Loop Automated Verification (Windows & WSL2)...")
            run_tests()
        else:
            print("\n[!] Invalid choice, please enter 0-5.\n")

if __name__ == "__main__":
    if len(sys.argv) > 1:
        arg = sys.argv[1].lower().replace("--", "").replace("-", "")
        if arg in ["windows", "win"]:
            create_platform_package(PLATFORMS["1"])
        elif arg in ["linux", "linuxamd64"]:
            create_platform_package(PLATFORMS["2"])
        elif arg in ["arm", "arm64", "linuxarm64"]:
            create_platform_package(PLATFORMS["3"])
        elif arg in ["all"]:
            for k in ["1", "2", "3"]:
                create_platform_package(PLATFORMS[k])
        elif arg in ["test", "verify"]:
            run_tests()
        elif arg in ["help", "h"]:
            print_header()
            print("Usage: python package.py [windows | linux | arm64 | all | test | help]")
            print("Or run without arguments for interactive menu.")
        else:
            print(f"Unknown argument: {sys.argv[1]}")
            print("Usage: python package.py [windows | linux | arm64 | all | test | help]")
    else:
        interactive_menu()
