import socket
import json
import time
import urllib.request
import threading

def get_cctl_token():
    req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/authorize/login',
        data=json.dumps({'username':'admin','password':'admin'}).encode('utf-8'),
        headers={'Content-Type':'application/json'}, method='POST')
    with urllib.request.urlopen(req) as r:
        return json.loads(r.read())['result']['token']

token = get_cctl_token()
headers = {'X-Access-Token': token, 'Content-Type': 'application/json'}

# 1. Connect TCP socket to CCTL port 18884
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('127.0.0.1', 18884))
print("[TCP] Connected to CCTL 18884")

# 2. Start listener thread to receive downlink INVOKE_FUNCTION and reply
def listener():
    f = s.makefile('r')
    while True:
        line = f.readline()
        if not line:
            break
        print(f"\n[TCP RECV Downlink]: {line.strip()}")
        try:
            msg = json.loads(line)
            if msg.get("messageType") == "INVOKE_FUNCTION":
                func_id = msg.get("functionId")
                msg_id = msg.get("messageId")
                dev_id = msg.get("deviceId")
                reply = {
                    "messageType": "INVOKE_FUNCTION_REPLY",
                    "deviceId": dev_id,
                    "messageId": msg_id,
                    "output": True,
                    "success": True
                }
                reply_bytes = (json.dumps(reply) + "\n").encode('utf-8')
                s.sendall(reply_bytes)
                print(f"[TCP SENT Reply]: {reply}")
        except Exception as e:
            print("Listener error:", e)

threading.Thread(target=listener, daemon=True).start()

# 3. Send initial telemetry to register device as ONLINE
now_ms = int(time.time() * 1000)
init_msg = {
    "deviceId": "pump-03",
    "properties": {
        "pressure": 1.95,
        "flow": 88.5,
        "vibration": 0.45,
        "runState": "running"
    },
    "timestamp": now_ms,
    "messageId": f"init_{now_ms}"
}
s.sendall((json.dumps(init_msg) + "\n").encode('utf-8'))
print("[TCP] Sent telemetry, device should be ONLINE")
time.sleep(1)

# 4. Check device detail state in CCTL
d_req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/device-instance/pump-03/detail', headers={'X-Access-Token': token})
with urllib.request.urlopen(d_req) as r:
    dev = json.loads(r.read())['result']
    print("Device state:", dev.get('state'))

# 5. Invoke startPump function via CCTL REST API!
print("\n[CCTL API] Invoking function 'startPump'...")
inv_req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/device-instance/pump-03/function/startPump',
    data=json.dumps({}).encode('utf-8'), headers=headers, method='POST')
try:
    with urllib.request.urlopen(inv_req, timeout=5) as r:
        res = json.loads(r.read())
        print("[CCTL API] Invoke startPump SUCCESS:", res)
except urllib.error.HTTPError as e:
    print("[CCTL API] Invoke startPump HTTPError:", e.code, e.read().decode())

# 6. Invoke setPressureTarget function via CCTL REST API!
print("\n[CCTL API] Invoking function 'setPressureTarget' with targetPressure=2.35...")
inv_req2 = urllib.request.Request('http://127.0.0.1:9100/jetlinks/device-instance/pump-03/function/setPressureTarget',
    data=json.dumps({"targetPressure": 2.35}).encode('utf-8'), headers=headers, method='POST')
try:
    with urllib.request.urlopen(inv_req2, timeout=5) as r:
        res = json.loads(r.read())
        print("[CCTL API] Invoke setPressureTarget SUCCESS:", res)
except urllib.error.HTTPError as e:
    print("[CCTL API] Invoke setPressureTarget HTTPError:", e.code, e.read().decode())

time.sleep(1)
s.close()
print("\nTest completed.")
