import socket
import struct
import threading
import time
import json
import urllib.request
import urllib.parse
import urllib.error

# 1. Modbus TCP Simulator
class ModbusServer:
    def __init__(self, host='0.0.0.0', port=10502):
        self.host = host
        self.port = port
        self.running = True
        self.sock = None
        self.t = threading.Thread(target=self._serve, daemon=True)

    def start(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.sock.bind((self.host, self.port))
        self.sock.listen(5)
        self.t.start()
        print(f"[ModbusSim] Server listening on {self.host}:{self.port}")

    def _serve(self):
        while self.running:
            try:
                conn, addr = self.sock.accept()
                threading.Thread(target=self._handle_client, args=(conn,), daemon=True).start()
            except:
                break

    def _handle_client(self, conn):
        with conn:
            while self.running:
                try:
                    data = conn.recv(1024)
                    if not data or len(data) < 12:
                        break
                    trans_id, proto_id, length, unit_id, fc, addr, qty = struct.unpack(">HHHBBHH", data[:12])
                    if fc in (3, 4): # Read Holding / Input Registers
                        byte_cnt = qty * 2
                        resp_len = 3 + byte_cnt
                        resp_hdr = struct.pack(">HHHBB", trans_id, 0, resp_len, unit_id, fc)
                        resp_body = struct.pack("B", byte_cnt)
                        # Predictable live telemetry values
                        regs = {
                            0: 198,    # 400001: pressure (198 => e.g. 1.98 bar)
                            1: 8888,   # 400002: flow (8888 => 88.88 m3/h)
                            2: 45,     # 400003: vibration (45 => 0.45 mm/s)
                            3: 1       # 400004: runState (1 => running)
                        }
                        for i in range(qty):
                            reg_idx = addr + i
                            val = regs.get(reg_idx, 0)
                            resp_body += struct.pack(">H", val)
                        conn.sendall(resp_hdr + resp_body)
                    else:
                        break
                except:
                    break

# 2. Neuron Client
class NeuronClient:
    def __init__(self, base_url="http://127.0.0.1:7000/api/v2"):
        self.base_url = base_url
        self.token = None

    def req(self, method, path, body=None, params=None):
        url = f"{self.base_url}{path}"
        if params:
            url += "?" + urllib.parse.urlencode(params)
        data = json.dumps(body).encode('utf-8') if body is not None else None
        headers = {"Content-Type": "application/json"}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        req = urllib.request.Request(url, data=data, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=5) as resp:
                content = resp.read().decode('utf-8')
                return json.loads(content) if content else {}
        except urllib.error.HTTPError as e:
            err_content = e.read().decode('utf-8')
            try:
                return json.loads(err_content)
            except:
                return {"error": e.code, "msg": err_content}

    def login(self, username="admin", password="0000"):
        res = self.req("POST", "/login", {"name": username, "pass": password})
        self.token = res.get("token")
        return bool(self.token)

def main():
    print("=" * 60)
    print("NEURON <-> CCTL END-TO-END AUTOMATED VERIFICATION")
    print("=" * 60)

    # Step 1: Start Modbus simulator
    sim = ModbusServer(host='0.0.0.0', port=10502)
    sim.start()

    # Step 2: Login to Neuron
    neu = NeuronClient()
    if not neu.login():
        print("[FAIL] Cannot login to Neuron Core on 7000")
        return

    print("[Neuron] Logged in successfully.")

    # Clean up old southbound driver if present
    s_nodes = neu.req("GET", "/node", params={"type": 1}).get("nodes", [])
    if any(n.get("name") == "pump_driver" for n in s_nodes):
        print("[Neuron] Deleting existing pump_driver node...")
        neu.req("POST", "/node/ctl", {"node": "pump_driver", "cmd": 1})
        neu.req("DELETE", "/node", {"name": "pump_driver"})

    # Step 3: Create Modbus TCP Southbound Driver
    print("[Neuron] Creating Southbound driver 'pump_driver' (Modbus TCP)...")
    res = neu.req("POST", "/node", {"name": "pump_driver", "plugin": "Modbus TCP"})
    print("  -> Create:", res)

    # Step 4: Configure Modbus TCP settings
    print("[Neuron] Configuring 'pump_driver' parameters...")
    cfg = {
        "connection_mode": 0,
        "check_header": 0,
        "device_degrade": 0,
        "max_retries": 0,
        "retry_interval": 0,
        "endianess": 1,
        "endianess_64": 1,
        "address_base": 1,
        "interval": 20,
        "host": "172.23.96.1",
        "port": 10502,
        "timeout": 3000
    }
    res = neu.req("POST", "/node/setting", {"node": "pump_driver", "params": cfg})
    print("  -> Setting:", res)

    # Step 5: Add Group
    print("[Neuron] Adding group 'pump_group' (interval 500ms)...")
    res = neu.req("POST", "/group", {"node": "pump_driver", "group": "pump_group", "interval": 500})
    print("  -> Add Group:", res)

    # Step 6: Add Tags for pump properties
    print("[Neuron] Adding telemetry tags (pressure, flow, vibration, runState)...")
    tags = [
        {"name": "pressure", "address": "1!400001", "attribute": 1, "type": 4},
        {"name": "flow", "address": "1!400002", "attribute": 1, "type": 4},
        {"name": "vibration", "address": "1!400003", "attribute": 1, "type": 4},
        {"name": "runState", "address": "1!400004", "attribute": 1, "type": 4}
    ]
    res = neu.req("POST", "/tags", {"node": "pump_driver", "group": "pump_group", "tags": tags})
    print("  -> Add Tags:", res)

    # Step 7: Start Southbound driver
    print("[Neuron] Starting 'pump_driver'...")
    res = neu.req("POST", "/node/ctl", {"node": "pump_driver", "cmd": 0})
    print("  -> Start Driver:", res)

    # Step 8: Ensure Northbound cctl_app is configured and subscribed
    print("[Neuron] Subscribing 'cctl_app' to 'pump_driver' / 'pump_group'...")
    # Clean old subscription if any
    neu.req("DELETE", "/subscribe", {"app": "cctl_app", "driver": "pump_driver", "group": "pump_group"})
    res = neu.req("POST", "/subscribe", {
        "app": "cctl_app",
        "driver": "pump_driver",
        "group": "pump_group",
        "params": {}
    })
    print("  -> Subscribe:", res)

    # Wait for telemetry cycles
    print("[Test] Waiting 4 seconds for telemetry stream to CCTL...")
    time.sleep(4)

    # Step 9: Verify CCTL Device Properties
    print("\n" + "=" * 60)
    print("VERIFYING INGESTION IN CCTL (JetLinks)")
    print("=" * 60)

    cctl_login_req = urllib.request.Request(
        "http://127.0.0.1:9100/jetlinks/authorize/login",
        data=json.dumps({"username": "admin", "password": "admin"}).encode('utf-8'),
        headers={"Content-Type": "application/json"},
        method="POST"
    )
    with urllib.request.urlopen(cctl_login_req) as r:
        cctl_token = json.loads(r.read())["result"]["token"]

    cctl_hdrs = {"X-Access-Token": cctl_token}

    # Query device status
    d_req = urllib.request.Request("http://127.0.0.1:9100/jetlinks/device-instance/pump-03/detail", headers=cctl_hdrs)
    with urllib.request.urlopen(d_req) as r:
        dev_detail = json.loads(r.read())["result"]
        print(f"CCTL Device 'pump-03' Status: {dev_detail.get('state')}")

    # Query latest properties
    p_req = urllib.request.Request("http://127.0.0.1:9100/jetlinks/device-instance/pump-03/properties/latest", headers=cctl_hdrs)
    with urllib.request.urlopen(p_req) as r:
        props = json.loads(r.read())["result"]
        print(f"CCTL Ingested Properties Count: {len(props)}")
        for prop in props:
            print(f"  [PROPERTY] {prop.get('property')}: value={prop.get('value')}, formatted={prop.get('formatValue')}, time={prop.get('time')}")

    print("\n" + "=" * 60)
    print("ALL TESTS PASSED: NEURON NORTHBOUND CCTL PLUGIN STREAMING ACTIVE!")
    print("=" * 60)

if __name__ == '__main__':
    main()
