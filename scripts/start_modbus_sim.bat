@echo off
chcp 65001 >nul
echo [10502] 正在启动 Modbus 多从机仿真器 (端口 10502)...
cd /d "%~dp0.."
start "Modbus Simulator (10502)" python "%~dp0..\modbus_sim.py"
echo 启动完成，从机地址: 1, 2, 3 监听端口: 10502
timeout /t 2 >nul
