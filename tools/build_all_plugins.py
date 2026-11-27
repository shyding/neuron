import os
import glob
import subprocess

ROOT_DIR = "/mnt/c/ss/learn/neuron"
PLUGINS_DIR = os.path.join(ROOT_DIR, "plugins")
BUILD_DIR = "/tmp/neuron_plugins_built"
os.makedirs(BUILD_DIR, exist_ok=True)

ENTRY_FILES = {
    "EtherNet-IP-1400.so": ("EtherNet_IP", ["enet_1400_plugin.c", "enet_1400_stack.c", "enet_1400_tag.c", "enet.c", "cip.c"]),
    "EtherNet-IP-5500.so": ("EtherNet_IP", ["enet_55_plugin.c", "enet_55_stack.c", "enet_55_tag.c", "enet.c", "cip.c"]),
    "EtherNet-IP.so": ("EtherNet_IP", ["enet_plugin.c", "enet_stack.c", "enet_tag.c", "enet.c", "cip.c"]),
    "HJ212.so": ("hj", ["hj_plugin.c", "hj.c", "hj_stack.c", "hj_tag.c"]),
    "a1e.so": ("a1e", ["a1e_plugin.c", "a1e_point.c", "a1e_protocol.c"]),
    "ab5000.so": ("cip", ["ab5000_plugin.c", "ab_stack.c", "cip.c", "enip.c", "tag.c"]),
    "ads.so": ("ads", ["ads_plugin.c", "ads.c", "ads_client.c", "ads_tag.c"]),
    "aws-iot.so": ("mqtt", ["aws_iot_plugin.c", "mqtt_config.c", "mqtt_handle.c", "mqtt_plugin_intf.c", "schema.c", "ptformat.pb-c.c"]),
    "azure-iot.so": ("mqtt", ["azure_iot_plugin.c", "mqtt_config.c", "mqtt_handle.c", "mqtt_plugin_intf.c", "schema.c", "ptformat.pb-c.c"]),
    "bacnet-scan.so": ("bacnet", ["scanning/scan_plugin.c", "bacapp.c", "bacnet.c", "bacnet_point.c", "bacnet_stack.c"]),
    "bacnet.so": ("bacnet", ["bacnet_plugin.c", "bacapp.c", "bacnet.c", "bacnet_point.c", "bacnet_stack.c"]),
    "brother.so": ("brother", ["brother_plugin.c", "brother.c", "brother_stack.c", "brother_tag.c"]),
    "cctl.so": ("cctl", ["cctl_plugin.c", "cctl_handle.c", "cctl_config.c", "cctl_tcp.c"]),
    "codesys3.so": ("codesys3", ["codesys3_plugin.c", "codesys3.c", "codesys3_point.c", "codesys3_stack.c"]),
    "comli.so": ("comli", ["comli_plugin.c", "comli.c", "comli_point.c"]),
    "df1.so": ("df1", ["df1_plugin.c", "df1.c", "df1_point.c", "df1_req.c", "df1_stack.c"]),
    "dlt645-1997.so": ("dlt645", ["dlt645_97.c", "dlt645.c", "dlt645_point.c", "dlt645_req.c", "dlt645_stack.c"]),
    "dlt645-2007.so": ("dlt645", ["dlt645_plugin.c", "dlt645.c", "dlt645_point.c", "dlt645_req.c", "dlt645_stack.c"]),
    "dnp3.so": ("dnp3", ["dnp3_plugin.c", "dnp3.c", "dnp3_point.c", "dnp3_stack.c"]),
    "fins-tcp.so": ("fins", ["fins_plugin.c", "fins.c", "fins_point.c", "fins_stack.c"]),
    "fins-udp.so": ("fins", ["fins_udp.c", "fins.c", "fins_point.c", "fins_stack.c"]),
    "focas.so": ("focas", ["focas_plugin.c", "focas_stack.c", "focas_tag.c"]),
    "fx.so": ("fx", ["fx_plugin.c", "fx_point.c", "fx_protocol.c"]),
    "gewu.so": ("gewu", ["gewu_plugin.c", "gewu_config.c", "gewu_handle.c"]),
    "heidenhain_cnc.so": ("heidenhain_cnc", ["heidenhain_cnc_plugin.c", "lsv2.c", "lsv2_point.c", "lsv2_stack.c"]),
    "hollysys-modbus-rtu.so": ("modbus", ["hollysys_modbus_rtu.c", "hollysys_modbus_point.c", "modbus.c", "modbus_point.c", "modbus_req.c", "modbus_stack.c", "xinje_modbus_point.c"]),
    "hollysys-modbus-tcp.so": ("modbus", ["hollysys_modbus_tcp.c", "hollysys_modbus_point.c", "modbus.c", "modbus_point.c", "modbus_req.c", "modbus_stack.c", "xinje_modbus_point.c"]),
    "hostlink-cmode.so": ("hostlink", ["hostlink_cmode.c", "hostlink.c", "hostlink_point.c", "hostlink_req.c", "hostlink_stack.c"]),
    "hsms.so": ("secs_gem", ["hsms_plugin.c", "hsms.c", "hsms_stack.c", "secs2.c", "secs_gem_point.c"]),
    "iec101.so": ("iec101", ["iec101_plugin.c", "iec101.c", "iec101_stack.c"]),
    "iec102.so": ("iec102", ["iec102_plugin.c", "iec102.c", "iec102_stack.c", "iec102_tag.c"]),
    "iec103.so": ("iec103", ["iec103_plugin.c", "iec103.c", "iec103_stack.c", "iec103_tag.c"]),
    "iec104-ms.so": ("iec104_ms", ["iec104_ms_plugin.c", "iec104_ms_point.c", "iec104_ms_stack.c", "../iec104/iec104.c"]),
    "iec104-standard.so": ("iec104", ["iec104_standard_plugin.c", "iec104.c", "iec104_stack.c"]),
    "iec61850.so": ("iec61850", ["iec61850_plugin.c", "asn.c", "iec61850_protocol.c", "iec61850_stack.c", "tpkt.c"]),
    "inovance-modbus-tcp.so": ("modbus", ["inovance_modbus_tcp.c", "hollysys_modbus_point.c", "modbus.c", "modbus_point.c", "modbus_req.c", "modbus_stack.c", "xinje_modbus_point.c"]),
    "knd.so": ("knd", ["knd_plugin.c", "knd_handle.c", "knd_point.c"]),
    "knx.so": ("knx", ["knx_plugin.c", "knx_cemi.c", "knx_client.c", "knx_frame.c"]),
    "kuka.so": ("kuka", ["kuka_plugin.c", "kuka_point.c", "kuka_stack.c"]),
    "mazak-udp.so": ("mazak_udp", ["mazak_udp_plugin.c", "mazak_udp_config.c", "mazak_udp_parser.c", "mazak_udp_tag.c"]),
    "mewtocol.so": ("mewtocol", ["mewtocol_plugin.c", "mewtocol_point.c", "mewtocol_protocol.c"]),
    "modbus-ascii.so": ("modbus", ["modbus_ascii.c", "hollysys_modbus_point.c", "modbus.c", "modbus_point.c", "modbus_req.c", "modbus_stack.c", "xinje_modbus_point.c"]),
    "modbus-qh-tcp.so": ("modbus", ["modbus_qh_tcp.c", "hollysys_modbus_point.c", "modbus.c", "modbus_point.c", "modbus_req.c", "modbus_stack.c", "xinje_modbus_point.c"]),
    "mpi.so": ("mpi", ["mpi_plugin.c", "mpi_point.c", "mpi_protocol.c"]),
    "mtconnect.so": ("mtconnect", ["mtc_plugin.c", "mtc_point.c", "mtc_stack.c"]),
    "neuhub.so": ("neuhub", ["neuhub_plugin.c", "neuhub.c", "neuhub_point.c", "neuhub_stack.c"]),
    "nona11.so": ("nona11", ["nona11_plugin.c", "nona11.c", "nona11_stack.c", "nona11_tag.c"]),
    "omron-cip.so": ("cip", ["omron_plugin.c", "cip.c", "enip.c", "stack.c", "tag.c"]),
    "opcua-server.so": ("opcua_server", ["opcua_server_plugin.c", "cert_util.c", "my_ua_accesscontrol.c", "opcua_server_handle.c", "opcua_server_node_manager.c"]),
    "opcua.so": ("opcua", ["opcua_plugin.c", "opcua_browse_cache.c", "opcua_client_common.c", "opcua_client_subscribe.c", "opcua_condition.c", "opcua_handle.c", "open62541_client.c"]),
    "qna3e.so": ("qna3e", ["qna3e_plugin.c", "qna3e_point.c", "qna3e_protocol.c"]),
    "qna4e.so": ("qna4e", ["qna4e_plugin.c", "qna4e_point.c", "qna4e_protocol.c"]),
    "s5fetch-write.so": ("s5fw", ["s5fw_plugin.c", "s5fw_point.c", "s5fw_protocol.c"]),
    "s7comm-for-300.so": ("s7comm", ["s7comm_300_plugin.c", "cotp.c", "s7comm.c", "s7comm_point.c", "s7comm_stack.c", "tpkt.c"]),
    "s7comm.so": ("s7comm", ["s7comm_plugin.c", "cotp.c", "s7comm.c", "s7comm_point.c", "s7comm_stack.c", "tpkt.c"]),
    "snmp.so": ("snmp", ["snmp_plugin.c", "asn.c", "snmp.c", "snmp_stack.c", "snmp_tag.c"]),
    "sparkplugb.so": ("sparkplugb", ["sparkplugb_plugin.c", "sparkplugb.c", "sparkplugb.pb.c", "util.c"]),
    "srtp.so": ("srtp", ["srtp_plugin.c", "srtp.c", "srtp_point.c", "srtp_stack.c"]),
    "websocket.so": ("websocket", ["websocket_plugin.c", "websocket_client.c", "websocket_config.c", "websocket_handle.c"]),
    "xinje-modbus-rtu.so": ("modbus", ["xinje_modbus_rtu.c", "hollysys_modbus_point.c", "modbus.c", "modbus_point.c", "modbus_req.c", "modbus_stack.c", "xinje_modbus_point.c"]),
}

print(f"Building {len(ENTRY_FILES)} industrial protocol plugins...")

inc_flags = [
    f"-I{ROOT_DIR}/include",
    f"-I{ROOT_DIR}/include/neuron",
    f"-I{ROOT_DIR}/plugins/include",
    f"-I{ROOT_DIR}/plugins",
]

results = {}
for so_name, (pdir, src_files) in sorted(ENTRY_FILES.items()):
    out_so = f"{BUILD_DIR}/libplugin-{so_name}"
    
    src_paths = []
    for f in src_files:
        p = os.path.join(PLUGINS_DIR, pdir, f) if not f.startswith("..") else os.path.join(PLUGINS_DIR, pdir, f)
        if os.path.exists(p):
            src_paths.append(p)
        else:
            # check root plugins
            alt = os.path.join(PLUGINS_DIR, f)
            if os.path.exists(alt):
                src_paths.append(alt)

    cmd = ["gcc", "-shared", "-fPIC"] + inc_flags + [f"-I{PLUGINS_DIR}/{pdir}"] + src_paths + [
        "-L/usr/local/lib", "-lneuron-base", "-o", out_so
    ]
    
    proc = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
    if proc.returncode == 0:
        results[so_name] = "OK"
    else:
        results[so_name] = f"FAIL: {proc.stderr[:100]}"

ok_count = sum(1 for v in results.values() if v == "OK")
fail_count = sum(1 for v in results.values() if v != "OK")

print(f"Build Results: {ok_count} OK, {fail_count} FAILED")
for so, r in sorted(results.items()):
    if r != "OK":
        print(f"  {so}: {r}")
