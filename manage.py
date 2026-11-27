#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Neuron & CCTL Multi-Environment Management Console
Provides centralized management for Ports 7000, 8085, 8088, 10502 and CCTL 9100/18884.
"""

import os
import sys
import time
import socket
import subprocess
import urllib.request
import urllib.error
import json

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
PYTHON_EXE = sys.executable

def check_port(port, host="127.0.0.1", timeout=0.8):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(timeout)
        res = s.connect_ex((host, port))
        s.close()
        return res == 0
    except:
        return False

def get_cctl_devices_status():
    if not check_port(9100):
        return {}
    try:
        login_url = "http://127.0.0.1:9100/jetlinks/authorize/login"
        req = urllib.request.Request(
            login_url,
            data=json.dumps({"username": "admin", "password": "admin"}).encode('utf-8'),
            headers={"Content-Type": "application/json"},
            method="POST"
        )
        with urllib.request.urlopen(req, timeout=2) as resp:
            token = json.loads(resp.read().decode('utf-8'))["result"]["token"]

        headers = {"X-Access-Token": token}
        devices = ["pump-02", "pump-03", "gw-pump-station", "sub-pump-01", "sub-pump-02", "sub-valve-01"]
        states = {}
        for dev in devices:
            try:
                r = urllib.request.Request(f"http://127.0.0.1:9100/jetlinks/device-instance/{dev}/detail", headers=headers)
                with urllib.request.urlopen(r, timeout=1.5) as res:
                    d = json.loads(res.read().decode('utf-8'))["result"]
                    states[dev] = d.get("state", {}).get("text", "离线")
            except:
                states[dev] = "未知"
        return states
    except:
        return {}

def start_7000():
    print("[7000] 正在启动 WSL Neuron C Core (端口 7000)...")
    if os.name == 'nt':
        subprocess.Popen(['wsl.exe', '-d', 'Ubuntu', 'bash', '-c', 'cd /opt/neuronex-4.0.2 && ./bin/neuron'], creationflags=subprocess.CREATE_NEW_CONSOLE)
    else:
        subprocess.Popen(['/opt/neuronex-4.0.2/bin/neuron'])
    time.sleep(2)
    if check_port(7000):
        print("  -> 7000 端口已就绪 (RUNNING)")
    else:
        print("  -> 7000 启动中，请稍后检查状态")

def stop_7000():
    print("[7000] 正在停止 WSL Neuron C Core...")
    cmd = 'wsl -d Ubuntu -e bash -c "pkill -9 neuron 2>/dev/null || true"'
    subprocess.run(cmd, shell=True)
    time.sleep(1)
    print("  -> 已发送停止指令")

def start_8085():
    print("[8085] 正在启动 WSL 官方 NeuronEX (端口 8085)...")
    if os.name == 'nt':
        subprocess.Popen(['wsl.exe', '-d', 'Ubuntu', 'bash', '-c', 'cd /opt/neuronex-4.0.2 && ./bin/neuronex alone'], creationflags=subprocess.CREATE_NEW_CONSOLE)
    else:
        subprocess.Popen(['/opt/neuronex-4.0.2/bin/neuronex', 'alone'])
    time.sleep(2)
    if check_port(8085):
        print("  -> 8085 端口已就绪 (RUNNING)")
    else:
        print("  -> 8085 启动中，请稍后检查状态")

def stop_8085():
    print("[8085] 正在停止 WSL 官方 NeuronEX...")
    cmd = 'wsl -d Ubuntu -e bash -c "pkill -9 neuronex 2>/dev/null || true"'
    subprocess.run(cmd, shell=True)
    time.sleep(1)
    print("  -> 已发送停止指令")

def start_8088():
    print("[8088] 正在启动 Windows 原生 Standalone NeuronEX (端口 8088)...")
    exe_path = os.path.join(BASE_DIR, "neuronex_windows_standalone", "bin", "neuronex.exe")
    cwd = os.path.join(BASE_DIR, "neuronex_windows_standalone")
    if not os.path.exists(exe_path):
        print(f"  -> 错误: 未找到 {exe_path}")
        return
    subprocess.Popen([exe_path, "alone"], cwd=cwd, creationflags=subprocess.CREATE_NEW_CONSOLE if os.name == 'nt' else 0)
    time.sleep(2)
    if check_port(8088):
        print("  -> 8088 端口已就绪 (RUNNING)")
    else:
        print("  -> 8088 启动中，请稍后检查状态")

def stop_8088():
    print("[8088] 正在停止 Windows 8088 原生服务...")
    if os.name == 'nt':
        subprocess.run('taskkill /F /IM neuronex.exe 2>nul', shell=True)
    print("  -> 已停止 8088 服务")

def start_10502():
    print("[10502] 正在启动 Modbus 多从机仿真器 (端口 10502)...")
    sim_script = os.path.join(BASE_DIR, "modbus_sim.py")
    subprocess.Popen([PYTHON_EXE, sim_script], cwd=BASE_DIR, creationflags=subprocess.CREATE_NEW_CONSOLE if os.name == 'nt' else 0)
    time.sleep(1.5)
    if check_port(10502):
        print("  -> 10502 端口已就绪 (RUNNING)")
    else:
        print("  -> 10502 启动中，请稍后检查状态")

def stop_10502():
    print("[10502] 正在停止 Modbus 仿真器...")
    # Find python process running modbus_sim.py
    if os.name == 'nt':
        cmd = 'powershell -Command "Get-CimInstance Win32_Process | Where-Object { $_.CommandLine -like \'*modbus_sim.py*\' } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force }"'
        subprocess.run(cmd, shell=True)
    print("  -> 已停止 10502 仿真器")

def start_all():
    print("\n========== 正在一键启动所有环境服务 ==========")
    start_10502()
    start_7000()
    start_8085()
    start_8088()
    print("================ 全部启动指令执行完毕 ================\n")

def stop_all():
    print("\n========== 正在一键停止所有服务 ==========")
    stop_8088()
    stop_8085()
    stop_7000()
    stop_10502()
    print("================ 全部服务已停止 ================\n")

def run_test():
    print("\n========== 运行全链路闭环自测试 ==========")
    test_script = os.path.join(BASE_DIR, "scratch", "test_all_downlinks.py")
    if os.path.exists(test_script):
        subprocess.run([PYTHON_EXE, test_script])
    else:
        print("未找到测试脚本，执行在线状态与属性检查...")
    telemetry_script = os.path.join(BASE_DIR, "scratch", "verify_all_telemetry.py")
    if os.path.exists(telemetry_script):
        subprocess.run([PYTHON_EXE, telemetry_script])
    print("================ 闭环测试结束 ================\n")

def print_dashboard():
    os.system('cls' if os.name == 'nt' else 'clear')
    p7000 = check_port(7000)
    p8085 = check_port(8085)
    p8088 = check_port(8088)
    p10502 = check_port(10502)
    p9100 = check_port(9100)
    p18884 = check_port(18884)

    def color_status(ok):
        return " [ 运行中 / OPEN ] " if ok else " [ 已停止 / CLOSED ]"

    print("=" * 72)
    print("     Neuron & CCTL 多环境集成管理控制台 (Management Console)     ")
    print("=" * 72)
    print(" 【核心环境与端口状态】:")
    print(f"   * 7000   Neuron C Core (WSL)           :{color_status(p7000)}")
    print(f"   * 8085   官方 NeuronEX 网关 (WSL)      :{color_status(p8085)}")
    print(f"   * 8088   Windows 原生 Standalone 网关  :{color_status(p8088)}")
    print(f"   * 10502  Modbus 多从机父子设备仿真器   :{color_status(p10502)}")
    print(f"   * 9100   CCTL (JetLinks) Web/REST 平台 :{color_status(p9100)}")
    print(f"   * 18884  CCTL TCP 协议网关入口         :{color_status(p18884)}")
    print("-" * 72)

    dev_states = get_cctl_devices_status()
    print(" 【CCTL 设备在线拓扑状态】:")
    if dev_states:
        for dev, st in dev_states.items():
            mark = "● 在线" if st == "在线" else "○ 离线"
            print(f"   * {dev:16} : {mark} ({st})")
    else:
        print("   (CCTL 平台未启动或未连接)")
    print("=" * 72)
    print(" 【管理菜单选项】:")
    print("   [1] 启动 7000 (Neuron C Core)        [2] 停止 7000")
    print("   [3] 启动 8085 (NeuronEX WSL)         [4] 停止 8085")
    print("   [5] 启动 8088 (Windows 原生)         [6] 停止 8088")
    print("   [7] 启动 10502 (Modbus 仿真器)       [8] 停止 10502")
    print("   --------------------------------------------------------")
    print("   [A] 一键启动全部环境 (10502 + 7000 + 8085 + 8088)")
    print("   [S] 一键停止全部环境")
    print("   [T] 运行全链路闭环测试 (遥测数据与下行控制指令测试)")
    print("   [R] 刷新状态面板")
    print("   [Q] 退出管理程序")
    print("=" * 72)

def main():
    while True:
        print_dashboard()
        try:
            choice = input("请输入操作编号: ").strip().upper()
        except (KeyboardInterrupt, EOFError):
            print("\n已退出。")
            break

        if choice == '1':
            start_7000()
        elif choice == '2':
            stop_7000()
        elif choice == '3':
            start_8085()
        elif choice == '4':
            stop_8085()
        elif choice == '5':
            start_8088()
        elif choice == '6':
            stop_8088()
        elif choice == '7':
            start_10502()
        elif choice == '8':
            stop_10502()
        elif choice == 'A':
            start_all()
        elif choice == 'S':
            stop_all()
        elif choice == 'T':
            run_test()
        elif choice == 'R':
            pass
        elif choice == 'Q':
            print("再见！")
            break
        else:
            print("无效输入，请重新输入。")
        
        if choice not in ['R', 'Q']:
            input("\n按回车键返回主菜单...")

if __name__ == '__main__':
    main()
