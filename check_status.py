import socket
import urllib.request
import json

def check_ports():
    ports = [7000, 8085, 8086, 8088, 9100, 18884]
    for p in ports:
        s = socket.socket()
        s.settimeout(0.5)
        res = s.connect_ex(('127.0.0.1', p))
        s.close()
        print(f"Port {p}: {'OPEN' if res == 0 else 'CLOSED'}")

def check_cctl_device():
    login_url = "http://127.0.0.1:9100/jetlinks/authorize/login"
    login_data = json.dumps({"username": "admin", "password": "admin", "expires": 3600000}).encode("utf-8")
    req = urllib.request.Request(login_url, data=login_data, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req) as resp:
        res = json.loads(resp.read().decode("utf-8"))
        token = res.get("result", {}).get("token")
        
    req2 = urllib.request.Request(
        "http://127.0.0.1:9100/jetlinks/device-instance/pump-03/detail",
        headers={"X-Access-Token": token}
    )
    with urllib.request.urlopen(req2) as resp:
        dev = json.loads(resp.read().decode("utf-8")).get("result", {})
        print("Device name:", dev.get("name"))
        print("Device state:", dev.get("state", {}).get("text"))
        print("Device metadata keys:", list(json.loads(dev.get("metadata", "{}")).keys()) if dev.get("metadata") else "None")
        print("Device deriveMetadata keys:", list(json.loads(dev.get("deriveMetadata", "{}")).keys()) if dev.get("deriveMetadata") else "None")

        # Let's check metadata contents
        raw_meta = dev.get("metadata") or dev.get("deriveMetadata") or "{}"
        meta = json.loads(raw_meta)
        print("Functions in metadata:", [f.get("id") for f in meta.get("functions", [])])
        print("Properties in metadata:", [p.get("id") for p in meta.get("properties", [])])
        print("Events in metadata:", [e.get("id") for e in meta.get("events", [])])

if __name__ == "__main__":
    check_ports()
    print("-" * 40)
    check_cctl_device()
