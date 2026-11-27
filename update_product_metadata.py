import urllib.request
import json

def get_token():
    req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/authorize/login',
        data=json.dumps({'username':'admin','password':'admin'}).encode('utf-8'),
        headers={'Content-Type':'application/json'}, method='POST')
    with urllib.request.urlopen(req) as r:
        return json.loads(r.read())['result']['token']

token = get_token()
headers = {'X-Access-Token': token, 'Content-Type': 'application/json'}

# Fetch water-pump product
req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/device-product/water-pump', headers=headers)
with urllib.request.urlopen(req) as r:
    prod = json.loads(r.read())['result']

meta = json.loads(prod['metadata'])

# Define Functions
meta['functions'] = [
    {
        "id": "startPump",
        "name": "启动水泵",
        "description": "远程启动供水泵运行",
        "inputs": [],
        "output": {"type": "boolean"}
    },
    {
        "id": "stopPump",
        "name": "停止水泵",
        "description": "远程停止供水泵运行",
        "inputs": [],
        "output": {"type": "boolean"}
    },
    {
        "id": "setPressureTarget",
        "name": "设定目标压力",
        "description": "远程设定供水泵目标出口压力(MPa)",
        "inputs": [
            {
                "id": "targetPressure",
                "name": "目标压力(MPa)",
                "valueType": {"type": "double"}
            }
        ],
        "output": {"type": "boolean"}
    },
    {
        "id": "resetAlarm",
        "name": "复位告警",
        "description": "复位设备异常告警标志",
        "inputs": [],
        "output": {"type": "boolean"}
    }
]

# Define Events
meta['events'] = [
    {
        "id": "overpressure_alarm",
        "name": "超压告警",
        "description": "水泵出口压力超过安全阈值报警",
        "valueType": {
            "type": "object",
            "properties": [
                {"id": "pressure", "name": "当前压力", "valueType": {"type": "double"}},
                {"id": "threshold", "name": "报警阈值", "valueType": {"type": "double"}},
                {"id": "level", "name": "报警级别", "valueType": {"type": "string"}}
            ]
        }
    },
    {
        "id": "vibration_warning",
        "name": "振动预警",
        "description": "水泵机组振动速度超出预警范围",
        "valueType": {
            "type": "object",
            "properties": [
                {"id": "vibration", "name": "当前振动值", "valueType": {"type": "double"}},
                {"id": "level", "name": "预警等级", "valueType": {"type": "string"}}
            ]
        }
    },
    {
        "id": "fault_event",
        "name": "设备故障停机",
        "description": "水泵检测到严重故障紧急停机",
        "valueType": {
            "type": "object",
            "properties": [
                {"id": "faultCode", "name": "故障代码", "valueType": {"type": "int"}},
                {"id": "message", "name": "故障描述", "valueType": {"type": "string"}}
            ]
        }
    }
]

# Define Tags
meta['tags'] = [
    {
        "id": "location",
        "name": "安装位置",
        "description": "设备物理安装位置",
        "valueType": {"type": "string"}
    },
    {
        "id": "model",
        "name": "设备型号",
        "description": "供水泵设备出厂型号",
        "valueType": {"type": "string"}
    }
]

prod['metadata'] = json.dumps(meta, ensure_ascii=False)

put_req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/device-product/water-pump',
    data=json.dumps(prod).encode('utf-8'), headers=headers, method='PUT')
with urllib.request.urlopen(put_req) as r:
    res = json.loads(r.read())
    print("Update product metadata result:", res)
