import urllib.request
import json

LOGIN_URL = "http://127.0.0.1:9100/jetlinks/authorize/login"
BASE_URL = "http://127.0.0.1:9100/jetlinks"

req = urllib.request.Request(
    LOGIN_URL,
    data=json.dumps({"username": "admin", "password": "admin"}).encode('utf-8'),
    headers={"Content-Type": "application/json"},
    method="POST"
)

with urllib.request.urlopen(req) as resp:
    res = json.loads(resp.read().decode('utf-8'))
    token = res["result"]["token"]

headers = {"X-Access-Token": token}

detail_req = urllib.request.Request(f"{BASE_URL}/device-instance/pump-03/detail", headers=headers)
with urllib.request.urlopen(detail_req) as resp:
    detail = json.loads(resp.read().decode('utf-8'))
    state = detail["result"].get("state", {})
    print("Device state:", state)

props_req = urllib.request.Request(f"{BASE_URL}/device-instance/pump-03/properties/latest", headers=headers)
with urllib.request.urlopen(props_req) as resp:
    props = json.loads(resp.read().decode('utf-8'))
    print("Device latest properties count:", len(props.result if hasattr(props, 'result') else props.get("result", [])))
    for p in props.get("result", []):
        print(f"  {p.get('property')}: value={p.get('value')} (formatted={p.get('formatValue')}) time={p.get('time')}")
