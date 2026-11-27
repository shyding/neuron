#!/usr/bin/env python3
import sys, os, time, subprocess, json, urllib.request, shutil, zipfile, tarfile

ROOT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DIST_DIR = os.path.join(ROOT_DIR, "dist")
VERSION = "v4.0.2"

def print_banner(msg):
    print("\n" + "=" * 70)
    print(f"  {msg}")
    print("=" * 70)

def test_http_endpoint(port, name, method="GET", path="/", data=None, token=None, expect_code=200):
    url = f"http://127.0.0.1:{port}{path}"
    headers = {}
    body = None
    if token:
        headers["Authorization"] = f"Bearer {token}"
    if data:
        headers["Content-Type"] = "application/json"
        body = json.dumps(data).encode("utf-8")
    
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
    try:
        resp = opener.open(req, timeout=5)
        res_data = resp.read().decode("utf-8", errors="ignore")
        if resp.status == expect_code:
            print(f"   [PASS] {name} ({method} {path}) -> HTTP {resp.status}")
            return True, res_data
        else:
            print(f"   [FAIL] {name} ({method} {path}) -> Expected {expect_code}, got {resp.status}")
            return False, res_data
    except urllib.error.HTTPError as e:
        res_data = e.read().decode("utf-8", errors="ignore")
        if e.code == expect_code:
            print(f"   [PASS] {name} ({method} {path}) -> HTTP {e.code}")
            return True, res_data
        print(f"   [FAIL] {name} ({method} {path}) -> HTTP {e.code}: {res_data[:100]}")
        return False, res_data
    except Exception as e:
        print(f"   [FAIL] {name} ({method} {path}) -> Error: {e}")
        return False, str(e)

def verify_instance(port):
    print(f"   Verifying instance running on port {port}...")
    # 1. Ping
    ok, _ = test_http_endpoint(port, "Ping API", "GET", "/api/ping")
    if not ok: return False
    
    # 2. Auth status
    ok, _ = test_http_endpoint(port, "Auth Check (Public)", "GET", "/api/auth")
    if not ok: return False

    # 3. Static dashboard
    ok, _ = test_http_endpoint(port, "Web Dashboard HTML", "GET", "/web/common/")
    if not ok: return False

    # 4. Login
    ok, login_res = test_http_endpoint(port, "Login API", "POST", "/api/login", data={"name": "admin", "pass": "0000"})
    if not ok: return False
    token = json.loads(login_res).get("token")
    if not token:
        print("   [FAIL] Missing token in login response")
        return False
    print(f"   [PASS] Token retrieved: {token[:20]}...")

    # 5. RBAC
    ok, rbac_res = test_http_endpoint(port, "RBAC Permissions", "GET", "/api/rbac", token=token)
    if not ok: return False

    # 6. License query
    ok, lic_res = test_http_endpoint(port, "License Status", "GET", "/api/license", token=token)
    if not ok: return False
    lic = json.loads(lic_res)
    print(f"   [PASS] Active License Type: {lic.get('licenseType')}, Tags: {lic.get('totalTagCount')}")

    # 7. Activation code test
    test_code = "NX-IVGVCICUMVRWQ3TPNRXWO2LFOMQES3TDFZ6DCMBQGAYDAMD4GI4TSOJNGAYS2MBR-64C1EB1B"
    ok, _ = test_http_endpoint(port, "Activation Code Update", "POST", "/api/license/activation_code", data={"activationCode": test_code}, token=token)
    if not ok: return False

    return True

def test_windows_package():
    print_banner("1/2: Testing Windows Package on Native Windows (Port 8088)")
    zip_name = f"neuronex-{VERSION}-windows-amd64.zip"
    zip_path = os.path.join(DIST_DIR, zip_name)
    if not os.path.exists(zip_path):
        print(f"[ERROR] Windows package not found at: {zip_path}")
        return False
    
    test_dir = os.path.join(ROOT_DIR, "dist", "tmp_test_win")
    if os.path.exists(test_dir):
        shutil.rmtree(test_dir, ignore_errors=True)
    os.makedirs(test_dir, exist_ok=True)

    print(f"-> Extracting {zip_name} to {test_dir}...")
    with zipfile.ZipFile(zip_path, 'r') as zf:
        zf.extractall(test_dir)

    # Modify port to 8088 in test config
    cfg_path = os.path.join(test_dir, "etc", "neuronex.yaml")
    if os.path.exists(cfg_path):
        with open(cfg_path, 'r', encoding='utf-8') as f:
            cfg = f.read()
        cfg = cfg.replace("port: 8085", "port: 8088").replace("port: 8086", "port: 8088")
        with open(cfg_path, 'w', encoding='utf-8') as f:
            f.write(cfg)

    exe_path = os.path.join(test_dir, "bin", "neuronex.exe")
    licgen_path = os.path.join(test_dir, "bin", "neuronex-licgen.exe")

    if not os.path.exists(exe_path):
        print(f"[FAIL] neuronex.exe not found in extracted package!")
        return False

    # Test standalone license generator
    print(f"-> Verifying standalone tool: {licgen_path}")
    gen_proc = subprocess.run([licgen_path, "-customer", "Windows-Test-Co", "-out", os.path.join(test_dir, "test.lic")],
                              capture_output=True, text=True)
    if gen_proc.returncode != 0:
        print(f"[FAIL] neuronex-licgen failed: {gen_proc.stderr}")
        return False
    print("   [PASS] neuronex-licgen generated test.lic successfully")

    # Start neuronex.exe on port 8088
    print("-> Starting neuronex.exe on port 8088...")
    proc = subprocess.Popen([exe_path, "alone"], cwd=test_dir, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    time.sleep(2.5)

    try:
        passed = verify_instance(8088)
        if passed:
            print("\n>>> Windows Package Verification: 100% PASSED! <<<")
            return True
        else:
            print("\n>>> Windows Package Verification: FAILED! <<<")
            return False
    finally:
        print("-> Cleaning up Windows test process...")
        proc.terminate()
        try:
            proc.wait(timeout=3)
        except:
            proc.kill()
        shutil.rmtree(test_dir, ignore_errors=True)

def test_wsl_package():
    print_banner("2/2: Testing Linux Package inside WSL2 (Port 8087)")
    tar_name = f"neuronex-{VERSION}-linux-amd64.tar.gz"
    tar_path = os.path.join(DIST_DIR, tar_name)
    if not os.path.exists(tar_path):
        print(f"[ERROR] Linux package not found at: {tar_path}")
        return False

    wsl_tar_path = tar_path.replace("\\", "/").replace("C:", "/mnt/c")
    wsl_script = f"""
    set -e
    TEST_DIR="/tmp/neuronex_linux_test"
    rm -rf "$TEST_DIR"
    mkdir -p "$TEST_DIR"
    tar -xzf "{wsl_tar_path}" -C "$TEST_DIR"
    chmod +x "$TEST_DIR/bin/"* || true
    
    # Change port to 8087
    sed -i 's/port: 8085/port: 8087/g' "$TEST_DIR/etc/neuronex.yaml"
    sed -i 's/port: 8086/port: 8087/g' "$TEST_DIR/etc/neuronex.yaml"

    # Test licgen
    echo "-> Testing neuronex-licgen inside WSL..."
    "$TEST_DIR/bin/neuronex-licgen" -customer "WSL-Test-Co" -out "$TEST_DIR/wsl.lic" > /dev/null
    
    echo "-> Starting neuronex alone on port 8087..."
    cd "$TEST_DIR"
    nohup ./bin/neuronex alone > "$TEST_DIR/run.log" 2>&1 &
    sleep 3
    """

    res = subprocess.run(["wsl.exe", "-d", "Ubuntu", "-e", "bash", "-c", wsl_script], capture_output=True, text=True, errors="replace")
    if res.returncode != 0:
        print(f"[FAIL] WSL initialization failed: {res.stderr}")
        return False

    try:
        # Run verification from inside WSL
        wsl_verify_script = f"""
        python3 -c '
import urllib.request, json
req = urllib.request.Request("http://127.0.0.1:8087/api/ping")
assert urllib.request.urlopen(req).status == 200
print("   [PASS] WSL Ping API -> 200")

req = urllib.request.Request("http://127.0.0.1:8087/api/auth")
assert urllib.request.urlopen(req).status == 200
print("   [PASS] WSL Auth Check -> 200")

req = urllib.request.Request("http://127.0.0.1:8087/api/login", data=b"{{\\"name\\":\\"admin\\",\\"pass\\":\\"0000\\"}}", headers={{"Content-Type": "application/json"}})
token = json.loads(urllib.request.urlopen(req).read())["token"]
print("   [PASS] WSL Login Token ->", token[:20])

req = urllib.request.Request("http://127.0.0.1:8087/api/rbac", headers={{"Authorization": f"Bearer {{token}}"}})
assert urllib.request.urlopen(req).status == 200
print("   [PASS] WSL RBAC Permissions -> 200")

req = urllib.request.Request("http://127.0.0.1:8087/api/license", headers={{"Authorization": f"Bearer {{token}}"}})
lic = json.loads(urllib.request.urlopen(req).read())
assert lic["licenseType"] == "COMMERCIAL"
print("   [PASS] WSL License -> COMMERCIAL, Tags:", lic["totalTagCount"])
'
        """
        verify_res = subprocess.run(["wsl.exe", "-d", "Ubuntu", "-e", "bash", "-c", wsl_verify_script], capture_output=True, text=True, errors="replace")
        print(verify_res.stdout)
        if verify_res.returncode == 0:
            print("\n>>> WSL Linux Package Verification: 100% PASSED! <<<")
            return True
        else:
            print(f"[FAIL] WSL Verification Error: {verify_res.stderr}")
            return False
    finally:
        cleanup_script = """
        pkill -f '/tmp/neuronex_linux_test/bin/neuronex' || true
        rm -rf /tmp/neuronex_linux_test
        """
        subprocess.run(["wsl.exe", "-d", "Ubuntu", "-e", "bash", "-c", cleanup_script])

def main():
    target = sys.argv[1] if len(sys.argv) > 1 else "all"
    win_ok = True
    wsl_ok = True

    if target in ["all", "windows"]:
        win_ok = test_windows_package()
    if target in ["all", "linux", "wsl"]:
        wsl_ok = test_wsl_package()

    print_banner("CLOSED-LOOP VERIFICATION SUMMARY")
    if target in ["all", "windows"]:
        print(f"  • Windows Native Package (Port 8088): {'[PASSED] ✅' if win_ok else '[FAILED] ❌'}")
    if target in ["all", "linux", "wsl"]:
        print(f"  • Linux WSL Package (Port 8087)     : {'[PASSED] ✅' if wsl_ok else '[FAILED] ❌'}")
    print("=" * 70)

    if not (win_ok and wsl_ok):
        sys.exit(1)

if __name__ == "__main__":
    main()
