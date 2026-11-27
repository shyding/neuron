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

req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/protocol/official-protocol', headers={'X-Access-Token': token})
with urllib.request.urlopen(req) as r:
    proto = json.loads(r.read())['result']

new_script = """
codec.decoder(function (context) {
    var message = context.getMessage();
    var data = JSON.parse(message.payloadAsString());
    if (data.messageType === "INVOKE_FUNCTION_REPLY") {
        return {
            "messageType": "INVOKE_FUNCTION_REPLY",
            "deviceId": data.deviceId,
            "messageId": data.messageId,
            "output": data.output,
            "success": data.success !== false
        };
    }
    if (data.messageType === "WRITE_PROPERTY_REPLY") {
        return {
            "messageType": "WRITE_PROPERTY_REPLY",
            "deviceId": data.deviceId,
            "messageId": data.messageId,
            "properties": data.properties,
            "success": data.success !== false
        };
    }
    if (data.event) {
        return {
            "messageType": "EVENT",
            "deviceId": data.deviceId,
            "event": data.event,
            "data": data.data,
            "timestamp": data.timestamp,
            "messageId": data.messageId
        };
    }
    return {
        "messageType": "REPORT_PROPERTY",
        "deviceId": data.deviceId,
        "properties": data.properties,
        "timestamp": data.timestamp,
        "messageId": data.messageId
    };
});
codec.encoder("READ_PROPERTY", function (context) {
    var message = context.getMessage();
    return JSON.stringify({
        "messageType": "READ_PROPERTY",
        "deviceId": message.getDeviceId(),
        "properties": message.getProperties(),
        "messageId": message.getMessageId()
    }) + "\\n";
});
codec.encoder("WRITE_PROPERTY", function (context) {
    var message = context.getMessage();
    return JSON.stringify({
        "messageType": "WRITE_PROPERTY",
        "deviceId": message.getDeviceId(),
        "properties": message.getProperties(),
        "messageId": message.getMessageId()
    }) + "\\n";
});
codec.encoder("INVOKE_FUNCTION", function (context) {
    var message = context.getMessage();
    return JSON.stringify({
        "messageType": "INVOKE_FUNCTION",
        "deviceId": message.getDeviceId(),
        "functionId": message.getFunctionId(),
        "inputs": message.getInputs(),
        "messageId": message.getMessageId()
    }) + "\\n";
});
"""

proto['configuration']['script'] = new_script

put_req = urllib.request.Request('http://127.0.0.1:9100/jetlinks/protocol/official-protocol',
    data=json.dumps(proto).encode('utf-8'), headers=headers, method='PUT')
try:
    with urllib.request.urlopen(put_req) as r:
        print('Update protocol result:', json.loads(r.read()))
except urllib.error.HTTPError as e:
    print('HTTPError:', e.code, e.read().decode())
