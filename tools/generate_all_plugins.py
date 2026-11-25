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

# Parse source files per plugin
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

print(f"Loaded {len(metadata_all)} plugin metadata entries and {len(so_files_map)} source maps.")

# Excluded base plugins already in repo
BASE_PLUGINS = {
    "modbus-tcp.so", "modbus-rtu.so", "mqtt.so", "mqtt-v2.so",
    "ekuiper.so", "file.so", "monitor.so", "datalayers.so", "kafka.so"
}

target_plugins = {}
for so_path, meta in metadata_all.items():
    so_name = os.path.basename(so_path).replace("libplugin-", "")
    if so_name in BASE_PLUGINS:
        continue
    target_plugins[so_name] = meta

print(f"Target industrial plugins to generate: {len(target_plugins)}")
