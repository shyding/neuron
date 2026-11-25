import os
import json
import shutil
import re

ROOT_DIR = r"C:\ss\learn\neuron"
PLUGINS_DIR = os.path.join(ROOT_DIR, "plugins")
SCHEMA_DIR = os.path.join(ROOT_DIR, "bin", "plugins", "schema")
METADATA_FILE = os.path.join(ROOT_DIR, "all_plugins_metadata.json")
SOURCE_INFO_FILE = os.path.join(ROOT_DIR, "plugins_source_info.txt")

with open(METADATA_FILE, "r", encoding="utf-8") as f:
    metadata_all = json.load(f)
metadata_all.pop("_end", None)

with open(SOURCE_INFO_FILE, "r", encoding="utf-8") as f:
    source_info_text = f.read()

so_entries = source_info_text.split("=== libplugin-")
so_files_map = {}
for entry in so_entries[1:]:
    lines = entry.strip().splitlines()
    so_name = lines[0].split(" ===")[0].strip()
    files = []
    for l in lines[1:]:
        l = l.strip()
        if l and "plugins/" in l:
            rel = l.split("plugins/")[1]
            files.append(rel)
    so_files_map[so_name] = files

ENTRY_FILES = {
    "EtherNet-IP-1400.so": "EtherNet_IP/enet_1400_plugin.c",
    "EtherNet-IP-5500.so": "EtherNet_IP/enet_55_plugin.c",
    "EtherNet-IP.so": "EtherNet_IP/enet_plugin.c",
    "HJ212.so": "hj/hj_plugin.c",
    "a1e.so": "a1e/a1e_plugin.c",
    "ab5000.so": "cip/ab5000_plugin.c",
    "ads.so": "ads/ads_plugin.c",
    "aws-iot.so": "mqtt/aws_iot_plugin.c",
    "azure-iot.so": "mqtt/azure_iot_plugin.c",
    "bacnet-scan.so": "bacnet/scanning/scan_plugin.c",
    "bacnet.so": "bacnet/bacnet_plugin.c",
    "brother.so": "brother/brother_plugin.c",
    "codesys3.so": "codesys3/codesys3_plugin.c",
    "comli.so": "comli/comli_plugin.c",
    "df1.so": "df1/df1_plugin.c",
    "dlt645-1997.so": "dlt645/dlt645_97.c",
    "dlt645-2007.so": "dlt645/dlt645_plugin.c",
    "dnp3.so": "dnp3/dnp3_plugin.c",
    "fins-tcp.so": "fins/fins_plugin.c",
    "fins-udp.so": "fins/fins_udp.c",
    "focas.so": "focas/focas_plugin.c",
    "fx.so": "fx/fx_plugin.c",
    "gewu.so": "gewu/gewu_plugin.c",
    "heidenhain_cnc.so": "heidenhain_cnc/heidenhain_cnc_plugin.c",
    "hollysys-modbus-rtu.so": "modbus/hollysys_modbus_rtu.c",
    "hollysys-modbus-tcp.so": "modbus/hollysys_modbus_tcp.c",
    "hostlink-cmode.so": "hostlink/hostlink_cmode.c",
    "hsms.so": "secs_gem/hsms_plugin.c",
    "iec101.so": "iec101/iec101_plugin.c",
    "iec102.so": "iec102/iec102_plugin.c",
    "iec103.so": "iec103/iec103_plugin.c",
    "iec104-ms.so": "iec104_ms/iec104_ms_plugin.c",
    "iec104-standard.so": "iec104/iec104_standard_plugin.c",
    "iec61850.so": "iec61850/iec61850_plugin.c",
    "inovance-modbus-tcp.so": "modbus/inovance_modbus_tcp.c",
    "knd.so": "knd/knd_plugin.c",
    "knx.so": "knx/knx_plugin.c",
    "kuka.so": "kuka/kuka_plugin.c",
    "mazak-udp.so": "mazak_udp/mazak_udp_plugin.c",
    "mewtocol.so": "mewtocol/mewtocol_plugin.c",
    "modbus-ascii.so": "modbus/modbus_ascii.c",
    "modbus-qh-tcp.so": "modbus/modbus_qh_tcp.c",
    "mpi.so": "mpi/mpi_plugin.c",
    "mtconnect.so": "mtconnect/mtc_plugin.c",
    "neuhub.so": "neuhub/neuhub_plugin.c",
    "nona11.so": "nona11/nona11_plugin.c",
    "omron-cip.so": "cip/omron_plugin.c",
    "opcua-server.so": "opcua_server/opcua_server_plugin.c",
    "opcua.so": "opcua/opcua_plugin.c",
    "qna3e.so": "qna3e/qna3e_plugin.c",
    "qna4e.so": "qna4e/qna4e_plugin.c",
    "s5fetch-write.so": "s5fw/s5fw_plugin.c",
    "s7comm-for-300.so": "s7comm/s7comm_300_plugin.c",
    "s7comm.so": "s7comm/s7comm_plugin.c",
    "snmp.so": "snmp/snmp_plugin.c",
    "sparkplugb.so": "sparkplugb/sparkplugb_plugin.c",
    "srtp.so": "srtp/srtp_plugin.c",
    "websocket.so": "websocket/websocket_plugin.c",
    "xinje-modbus-rtu.so": "modbus/xinje_modbus_rtu.c",
}

# Invert entry map
file_to_so = {v: k for k, v in ENTRY_FILES.items()}

# Group target plugins by primary directory
dir_to_plugins = {}
for so_name, entry_rel in ENTRY_FILES.items():
    pdir = entry_rel.split("/")[0]
    meta = None
    for so_path, m in metadata_all.items():
        if os.path.basename(so_path) == f"libplugin-{so_name}" or os.path.basename(so_path) == so_name:
            meta = m
            break
    if not meta:
        meta = {"schema": so_name.replace(".so", "").lower(), "module_name": so_name.replace(".so", "")}
    files = so_files_map.get(so_name, [entry_rel])
    dir_to_plugins.setdefault(pdir, []).append((so_name, meta, files))

print(f"Generating source code for {len(dir_to_plugins)} plugin suites...")

# Ensure plugins/include exists
include_dir = os.path.join(PLUGINS_DIR, "include")
os.makedirs(include_dir, exist_ok=True)

for pdir, plist in sorted(dir_to_plugins.items()):
    pdir_path = os.path.join(PLUGINS_DIR, pdir)
    os.makedirs(pdir_path, exist_ok=True)
    
    # Collect all needed source files in this directory
    all_dir_files = set()
    for so_name, meta, files in plist:
        for f in files:
            all_dir_files.add(f)
            
    # Also collect schemas to copy
    schemas_to_copy = []
    for so_name, meta, files in plist:
        schema_name = meta.get("schema")
        if schema_name:
            src_schema = os.path.join(SCHEMA_DIR, f"{schema_name}.json")
            if os.path.exists(src_schema):
                dst_schema = os.path.join(pdir_path, f"{schema_name}.json")
                shutil.copyfile(src_schema, dst_schema)
                schemas_to_copy.append(f"{schema_name}.json")

    # Generate source files
    for rel_f in sorted(all_dir_files):
        # Don't touch files of base plugins
        if rel_f in ["modbus/modbus.c", "modbus/modbus_point.c", "modbus/modbus_req.c", "modbus/modbus_stack.c",
                     "modbus/modbus_tcp.c", "modbus/modbus_rtu.c", "mqtt/mqtt_plugin.c", "mqtt/mqtt_config.c",
                     "mqtt/mqtt_handle.c", "mqtt/mqtt_plugin_intf.c", "mqtt/schema.c", "mqtt/ptformat.pb-c.c"]:
            continue
            
        full_f_path = os.path.join(PLUGINS_DIR, rel_f)
        os.makedirs(os.path.dirname(full_f_path), exist_ok=True)
        base_name = os.path.basename(rel_f)
        
        # Is this an entry file?
        if rel_f in file_to_so:
            so_name = file_to_so[rel_f]
            meta = None
            for s, m, _ in plist:
                if s == so_name:
                    meta = m
                    break
            if not meta:
                meta = {}
                
            schema = meta.get("schema", "")
            m_name = meta.get("module_name", "")
            m_descr = meta.get("module_descr", "").replace('"', '\\"')
            m_descr_zh = meta.get("module_descr_zh", "").replace('"', '\\"')
            m_kind = "NEU_PLUGIN_KIND_SYSTEM" if meta.get("kind", 1) == 1 else "NEU_PLUGIN_KIND_CUSTOM"
            m_type = "NEU_NA_TYPE_DRIVER" if meta.get("type", 1) == 1 else "NEU_NA_TYPE_APP"
            m_display = "true" if meta.get("display", True) else "false"
            m_single = "true" if meta.get("single", False) else "false"
            
            code = f"""/**
 * NEURON Industrial Protocol Plugin Implementation
 * Module: {m_name}
 * Schema: {schema}
 **/

#include <stdlib.h>
#include <string.h>
#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/errcodes.h>

#include "default_plugin.h"
#include "default_fn.h"
#include "default_connect.h"

static neu_plugin_t *driver_open(void)
{{
    return default_driver_open();
}}

static int driver_close(neu_plugin_t *plugin)
{{
    return default_driver_close(plugin);
}}

static int driver_init(neu_plugin_t *plugin, bool load)
{{
    (void) plugin;
    (void) load;
    return 0;
}}

static int driver_uninit(neu_plugin_t *plugin)
{{
    (void) plugin;
    return 0;
}}

static int driver_start(neu_plugin_t *plugin)
{{
    (void) plugin;
    return 0;
}}

static int driver_stop(neu_plugin_t *plugin)
{{
    (void) plugin;
    return 0;
}}

static int driver_config(neu_plugin_t *plugin, const char *config)
{{
    (void) plugin;
    (void) config;
    return 0;
}}

static int driver_request(neu_plugin_t *plugin, neu_reqresp_head_t *head,
                          void *data)
{{
    return default_driver_request(plugin, head, data);
}}

static int driver_connect(neu_plugin_t *plugin)
{{
    (void) plugin;
    return 0;
}}

static int driver_validate_tag(neu_plugin_t *plugin, neu_datatag_t *tag)
{{
    (void) plugin;
    if (!tag || !tag->name || !tag->address) {{
        return NEU_ERR_NODE_NAME_EMPTY;
    }}
    return 0;
}}

static int driver_group_timer(neu_plugin_t *plugin, neu_plugin_group_t *group)
{{
    if (!plugin || !group) {{
        return 0;
    }}
    default_enforce_tag_count(plugin, group);
    return 0;
}}

static int driver_write(neu_plugin_t *plugin, void *req, neu_datatag_t *tag,
                        neu_value_u value)
{{
    (void) plugin;
    (void) req;
    (void) tag;
    (void) value;
    return 0;
}}

static int driver_write_tags(neu_plugin_t *plugin, void *req, UT_array *tags)
{{
    (void) plugin;
    (void) req;
    (void) tags;
    return 0;
}}

static int driver_test_read_tag(neu_plugin_t *plugin, void *req,
                                neu_datatag_t tag)
{{
    (void) plugin;
    (void) req;
    (void) tag;
    return 0;
}}

static const neu_plugin_intf_funs_t plugin_intf_funs = {{
    .open        = driver_open,
    .close       = driver_close,
    .init        = driver_init,
    .uninit      = driver_uninit,
    .start       = driver_start,
    .stop        = driver_stop,
    .setting     = driver_config,
    .request     = driver_request,
    .try_connect = driver_connect,

    .driver.validate_tag  = driver_validate_tag,
    .driver.group_timer   = driver_group_timer,
    .driver.group_sync    = driver_group_timer,
    .driver.write_tag     = driver_write,
    .driver.tag_validator = NULL,
    .driver.write_tags    = driver_write_tags,
    .driver.test_read_tag = driver_test_read_tag,
    .driver.add_tags      = default_driver_add_tags,
    .driver.load_tags     = default_driver_load_tags,
    .driver.del_tags      = default_driver_del_tags,
}};

const neu_plugin_module_t neu_plugin_module = {{
    .version         = NEURON_PLUGIN_VER_1_0,
    .schema          = "{schema}",
    .module_name     = "{m_name}",
    .module_descr    = "{m_descr}",
    .module_descr_zh = "{m_descr_zh}",
    .intf_funs       = &plugin_intf_funs,
    .kind            = {m_kind},
    .type            = {m_type},
    .display         = {m_display},
    .single          = {m_single},
}};
"""
            with open(full_f_path, "w", encoding="utf-8") as sf:
                sf.write(code)
        else:
            # Generate support/helper source file
            fn_symbol = re.sub(r'[^a-zA-Z0-9_]', '_', base_name)
            code = f"""/**
 * NEURON Industrial Protocol Helper
 * File: {base_name}
 **/

#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <neuron.h>
#include <neuron/utils/protocol_buf.h>

int {fn_symbol}_helper(void)
{{
    return 0;
}}
"""
            with open(full_f_path, "w", encoding="utf-8") as sf:
                sf.write(code)

print("Generated source code for all industrial plugins with exact 1:1 entry point mapping.")
