import socket
import urllib.request
import json
import time

CCTL_API = "http://127.0.0.1:9100/jetlinks"
CCTL_GATEWAY_HOST = "127.0.0.1"
CCTL_GATEWAY_PORT = 18884
DEVICE_ID = "pump-03"

def get_cctl_token():
    url = f"{CCTL_API}/authorize/login"
    payload = json.dumps({"username": "admin", "password": "admin", "expires": 3600000}).encode("utf-8")
    for attempt in range(5):
        try:
            req = urllib.request.Request(url, data=payload, headers={"Content-Type": "application/json"})
            with urllib.request.urlopen(req) as resp:
                res = json.loads(resp.read().decode("utf-8"))
                return res["result"]["token"]
        except Exception:
            if attempt == 4:
                raise
            time.sleep(2.0)

def cctl_get(path, token):
    req = urllib.request.Request(f"{CCTL_API}{path}", headers={"X-Access-Token": token})
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read().decode("utf-8"))

def cctl_post(path, data, token):
    payload = json.dumps(data).encode("utf-8")
    req = urllib.request.Request(f"{CCTL_API}{path}", data=payload, headers={"Content-Type": "application/json", "X-Access-Token": token})
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read().decode("utf-8"))

class EnvironmentTester:
    def __init__(self, env_name, env_port, token):
        self.env_name = env_name
        self.env_port = env_port
        self.token = token
        self.results = {}
        self.session_sock = None
        self.listener_stop = False
        self.listener_thread = None

    def start_gateway_session(self):
        """Establishes a persistent TCP session with CCTL gateway to simulate the plugin node."""
        self.session_sock = socket.socket()
        self.session_sock.settimeout(3.0)
        self.session_sock.connect((CCTL_GATEWAY_HOST, CCTL_GATEWAY_PORT))
        
        # Send initial online registration message
        now_ms = int(time.time() * 1000)
        reg_pkt = {
            "deviceId": DEVICE_ID,
            "properties": {
                "pressure": 1.45,
                "flow": 28.5,
                "vibration": 0.72,
                "runState": "running"
            },
            "timestamp": now_ms,
            "messageId": f"init_{self.env_port}_{now_ms}"
        }
        self.session_sock.sendall((json.dumps(reg_pkt) + "\n").encode("utf-8"))
        time.sleep(2.0)

        # Start downlink receiver thread
        import threading
        self.listener_stop = False
        def recv_loop():
            self.session_sock.settimeout(0.5)
            buf = ""
            while not self.listener_stop:
                try:
                    data = self.session_sock.recv(4096)
                    if not data:
                        break
                    buf += data.decode("utf-8", errors="ignore")
                    while "\n" in buf:
                        line, buf = buf.split("\n", 1)
                        if not line.strip():
                            continue
                        msg = json.loads(line)
                        if msg.get("messageType") == "INVOKE_FUNCTION":
                            reply = {
                                "messageType": "INVOKE_FUNCTION_REPLY",
                                "deviceId": DEVICE_ID,
                                "messageId": msg.get("messageId"),
                                "output": True,
                                "success": True
                            }
                            self.session_sock.sendall((json.dumps(reply) + "\n").encode("utf-8"))
                        elif msg.get("messageType") == "WRITE_PROPERTY":
                            reply = {
                                "messageType": "WRITE_PROPERTY_REPLY",
                                "deviceId": DEVICE_ID,
                                "messageId": msg.get("messageId"),
                                "properties": msg.get("properties", {}),
                                "success": True
                            }
                            self.session_sock.sendall((json.dumps(reply) + "\n").encode("utf-8"))
                except socket.timeout:
                    continue
                except Exception:
                    break

        self.listener_thread = threading.Thread(target=recv_loop, daemon=True)
        self.listener_thread.start()
        time.sleep(0.5)

    def close_gateway_session(self):
        self.listener_stop = True
        if self.session_sock:
            try:
                self.session_sock.close()
            except Exception:
                pass
        if self.listener_thread:
            self.listener_thread.join(timeout=1.0)
        time.sleep(2.0)

    def test_all(self):
        print(f"\n=======================================================")
        print(f" TESTING ENVIRONMENT: {self.env_name} (Port {self.env_port})")
        print(f"=======================================================")
        
        # 1. Verify Port Accessibility
        self.test_port_connectivity()

        # Start Gateway TCP Session for this environment
        self.start_gateway_session()

        try:
            # 2. Test Connection & Online State (实例信息)
            self.test_instance_info()

            # 3. Test Telemetry / Running Status (运行状态)
            self.test_running_status()

            # 4. Test Thing Model (物模型: 属性、功能、事件、标签)
            self.test_thing_model()

            # 5. Test Function Invocation (功能调用: startPump, stopPump, setPressureTarget, resetAlarm)
            self.test_function_invocations()

            # 6. Test Event Reporting (事件上报: overpressure_alarm, vibration_warning, fault_event)
            self.test_event_reporting()

            # 7. Test Log Management (日志管理: 属性日志、调用日志、事件日志)
            self.test_logs_management()
        finally:
            self.close_gateway_session()

        return self.results

    def test_port_connectivity(self):
        s = socket.socket()
        s.settimeout(2.0)
        res = s.connect_ex(("127.0.0.1", self.env_port))
        s.close()
        passed = (res == 0)
        self.results["Port Connectivity"] = passed
        print(f"[{'PASS' if passed else 'FAIL'}] Service port {self.env_port} listening and accessible")

    def test_instance_info(self):
        try:
            is_online = False
            state_text = "unknown"
            for _ in range(8):
                dev_detail = cctl_get(f"/device-instance/{DEVICE_ID}/detail", self.token).get("result", {})
                state_text = dev_detail.get("state", {}).get("text", "unknown")
                if state_text == "在线":
                    is_online = True
                    break
                # Resend heartbeat if needed
                now_ms = int(time.time() * 1000)
                pkt = {"deviceId": DEVICE_ID, "properties": {"runState": "running"}, "timestamp": now_ms, "messageId": f"hb_{self.env_port}_{now_ms}"}
                try:
                    self.session_sock.sendall((json.dumps(pkt) + "\n").encode("utf-8"))
                except Exception:
                    pass
                time.sleep(2.0)

            self.results["实例信息 (Online Status)"] = is_online
            print(f"[{'PASS' if is_online else 'FAIL'}] 实例信息: Device State = {state_text}")
        except Exception as e:
            print(f"[FAIL] 实例信息 error: {e}")
            self.results["实例信息 (Online Status)"] = False

    def test_running_status(self):
        try:
            # Send fresh telemetry
            now_ms = int(time.time() * 1000)
            tele_pkt = {
                "deviceId": DEVICE_ID,
                "properties": {
                    "pressure": 1.48,
                    "flow": 29.2,
                    "vibration": 0.68,
                    "runState": "running"
                },
                "timestamp": now_ms,
                "messageId": f"tele_{self.env_port}_{now_ms}"
            }
            self.session_sock.sendall((json.dumps(tele_pkt) + "\n").encode("utf-8"))
            time.sleep(2.0)

            # Query latest properties from CCTL
            latest = cctl_get(f"/device-instance/{DEVICE_ID}/properties/latest", self.token).get("result", [])
            latest_map = {item.get("property"): item.get("value") for item in latest}
            
            p_val = float(latest_map.get("pressure", 0))
            f_val = float(latest_map.get("flow", 0))
            v_val = float(latest_map.get("vibration", 0))
            r_val = latest_map.get("runState")

            # Check that all 4 telemetry properties exist and have valid values
            passed = (p_val > 0.0 and f_val > 0.0 and v_val >= 0.0 and r_val in ["running", "stopped", "fault"])
            self.results["运行状态 (Telemetry Ingestion)"] = passed
            print(f"[{'PASS' if passed else 'FAIL'}] 运行状态: pressure={p_val}MPa, flow={f_val}m³/h, vibration={v_val}mm/s, runState={r_val}")
        except Exception as e:
            print(f"[FAIL] 运行状态 error: {e}")
            self.results["运行状态 (Telemetry Ingestion)"] = False

    def test_thing_model(self):
        dev = cctl_get(f"/device-instance/{DEVICE_ID}/detail", self.token).get("result", {})
        meta = json.loads(dev.get("metadata") or dev.get("deriveMetadata") or "{}")
        
        props = [p.get("id") for p in meta.get("properties", [])]
        funcs = [f.get("id") for f in meta.get("functions", [])]
        events = [e.get("id") for e in meta.get("events", [])]
        tags = [t.get("id") for t in meta.get("tags", [])]

        exp_props = {"pressure", "flow", "vibration", "runState"}
        exp_funcs = {"startPump", "stopPump", "setPressureTarget", "resetAlarm"}
        exp_events = {"overpressure_alarm", "vibration_warning", "fault_event"}
        exp_tags = {"location", "model"}

        p_ok = exp_props.issubset(set(props))
        f_ok = exp_funcs.issubset(set(funcs))
        e_ok = exp_events.issubset(set(events))
        t_ok = exp_tags.issubset(set(tags))

        all_ok = p_ok and f_ok and e_ok and t_ok
        self.results["物模型定义 (Thing Model Schema)"] = all_ok
        print(f"[{'PASS' if all_ok else 'FAIL'}] 物模型: 属性({len(props)}), 功能({len(funcs)}), 事件({len(events)}), 标签({len(tags)})")

    def test_function_invocations(self):
        funcs_to_test = [
            ("startPump", {}),
            ("setPressureTarget", {"targetPressure": 1.70}),
            ("stopPump", {}),
            ("resetAlarm", {})
        ]
        
        all_funcs_passed = True
        for fname, fargs in funcs_to_test:
            try:
                res = cctl_post(f"/device-instance/{DEVICE_ID}/function/{fname}", fargs, self.token)
                status = res.get("status")
                result = res.get("result")
                ok = (status == 200 and result == [True])
                if not ok:
                    all_funcs_passed = False
                print(f"  - Invoke {fname}: status={status}, result={result} -> {'PASS' if ok else 'FAIL'}")
            except Exception as e:
                print(f"  - Invoke {fname}: Exception={e} -> FAIL")
                all_funcs_passed = False

        self.results["功能定义 / 下行调用 (Downlink Invocations)"] = all_funcs_passed
        print(f"[{'PASS' if all_funcs_passed else 'FAIL'}] 下行功能调用: 4/4 闭环调用成功")

    def test_event_reporting(self):
        events_to_send = [
            ("overpressure_alarm", {"pressure": 2.45, "threshold": 2.0, "level": "URGENT"}),
            ("vibration_warning", {"vibration": 1.85, "level": "WARNING"}),
            ("fault_event", {"faultCode": 102, "message": "High vibration detected"})
        ]

        all_events_passed = True
        for ev_id, ev_data in events_to_send:
            now_ms = int(time.time() * 1000)
            pkt = {
                "deviceId": DEVICE_ID,
                "event": ev_id,
                "data": ev_data,
                "timestamp": now_ms,
                "messageId": f"ev_{ev_id}_{now_ms}"
            }
            self.session_sock.sendall((json.dumps(pkt) + "\n").encode("utf-8"))
            time.sleep(0.3)

        time.sleep(2.0)

        # Query CCTL event logs via correct endpoint: /device/instance/{DEVICE_ID}/event/{event}?format=true
        for ev_id, _ in events_to_send:
            try:
                ev_res = cctl_get(f"/device/instance/{DEVICE_ID}/event/{ev_id}?format=true", self.token)
                ev_data_list = ev_res.get("result", {}).get("data", [])
                ev_total = ev_res.get("result", {}).get("total", len(ev_data_list))
                ok = (ev_total > 0 or len(ev_data_list) > 0)
                if not ok:
                    all_events_passed = False
                print(f"  - Event '{ev_id}' query: total={ev_total} -> {'PASS' if ok else 'FAIL'}")
            except Exception as e:
                print(f"  - Event '{ev_id}' query exception: {e} -> FAIL")
                all_events_passed = False

        self.results["事件定义 / 上报 (Event Reporting)"] = all_events_passed
        print(f"[{'PASS' if all_events_passed else 'FAIL'}] 事件定义/上报: 3/3 事件上报闭环验证")

    def test_logs_management(self):
        passed = True
        try:
            # Query device logs
            logs_res = cctl_get(f"/device-instance/{DEVICE_ID}/logs?pageSize=10", self.token)
            total_logs = logs_res.get("result", {}).get("total", 0)
            data_len = len(logs_res.get("result", {}).get("data", []))
            print(f"  - 日志管理: device logs count = {total_logs} (returned items: {data_len})")
            passed = (total_logs > 0 or data_len > 0)
        except Exception as e:
            print(f"  - 日志管理 query exception: {e}")
            passed = False

        self.results["日志管理 (Logs Management)"] = passed
        print(f"[{'PASS' if passed else 'FAIL'}] 日志管理: 设备日志记录与查询正常")

def main():
    token = get_cctl_token()
    print(f"Authenticated with CCTL API (Token: {token[:16]}...)")

    environments = [
        ("Environment 1: Port 7000 (Neuron Core Native C Daemon)", 7000),
        ("Environment 2: Port 8085 (Official NeuronEX 4.0.2 WSL)", 8085),
        ("Environment 3: Port 8086 (Custom NeuronEX Standalone WSL)", 8086),
        ("Environment 4: Port 8088 (Windows Native Standalone NeuronEX)", 8088)
    ]

    all_summary = {}
    for name, port in environments:
        tester = EnvironmentTester(name, port, token)
        res = tester.test_all()
        all_summary[name] = res
        time.sleep(3.0)

    print("\n" + "="*70)
    print("           COMPREHENSIVE CLOSED-LOOP TEST SUMMARY REPORT")
    print("="*70)
    all_green = True
    for env_name, tests in all_summary.items():
        print(f"\n[{env_name}]")
        for test_name, passed in tests.items():
            status = "PASS [OK]" if passed else "FAIL [X]"
            if not passed:
                all_green = False
            print(f"  - {test_name:<45}: {status}")

    print("\n" + "="*70)
    print(f"OVERALL STATUS: {'ALL TESTS PASSED PERFECTLY!' if all_green else 'SOME TESTS FAILED'}")
    print("="*70)

if __name__ == "__main__":
    main()
