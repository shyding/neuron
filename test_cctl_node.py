import json
import urllib.request
import urllib.parse
import urllib.error
import time

BASE_URL = "http://127.0.0.1:7000/api/v2"
token = None

def req(method, path, body=None, params=None):
    global token
    url = f"{BASE_URL}{path}"
    if params:
        url += "?" + urllib.parse.urlencode(params)
    data = json.dumps(body).encode('utf-8') if body is not None else None
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=5) as resp:
            content = resp.read().decode('utf-8')
            return json.loads(content) if content else {}
    except urllib.error.HTTPError as e:
        err_content = e.read().decode('utf-8')
        print(f"HTTPError {e.code} on {method} {url}: {err_content}")
        try:
            return json.loads(err_content)
        except:
            return {"error": e.code, "msg": err_content}
    except Exception as e:
        print(f"Exception on {method} {url}: {e}")
        return {"error": -1, "msg": str(e)}

print("=== 0. Logging in ===")
login_res = req("POST", "/login", {"name": "admin", "pass": "0000"})
token = login_res.get("token")
print("Token acquired:", bool(token))

print("\n=== 1. Checking existing nodes ===")
nodes = req("GET", "/node", params={"type": 2})
print("Existing nodes:", [n.get("name") for n in nodes.get("nodes", [])])

print("\n=== 4. Starting CCTL node ===")
res_start = req("POST", "/node/ctl", {"node": "cctl_app", "cmd": 0})
print("Start response:", res_start)

time.sleep(2)
print("\n=== 5. Checking node state ===")
node_state = req("GET", "/node/state", params={"node": "cctl_app"})
print("CCTL node state:", node_state)
