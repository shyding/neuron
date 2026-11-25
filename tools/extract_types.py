import gdb
import json
import re
import sys
import os

info_types = gdb.execute("info types", to_string=True)
info_funcs = gdb.execute("info functions", to_string=True)
info_vars = gdb.execute("info variables", to_string=True)

def parse_info(output):
    result = {}
    current_file = None
    for line in output.splitlines():
        if line.startswith("File "):
            current_file = line[5:].rstrip(":")
            if current_file not in result:
                result[current_file] = []
        elif current_file and line.strip() and not line.startswith("All defined ") and not line.startswith("Non-debugging "):
            result[current_file].append(line.strip())
    return result

types_by_file = parse_info(info_types)
funcs_by_file = parse_info(info_funcs)

target_files = [f for f in types_by_file if "neuron-modules" in f or "plugins" in f]
print(f"Target files: {len(target_files)}")

for f in target_files:
    print(f"\n==================== {f} ====================")
    print("TYPES:")
    for t in types_by_file.get(f, []):
        print(f"  {t}")
    print("FUNCTIONS:")
    for fn in funcs_by_file.get(f, []):
        print(f"  {fn}")
