@echo off
chcp 65001 >nul
echo [10502] 正在启动 Modbus 多从机仿真器 (端口 10502)...
cd /d "%~dp0.."
start "Modbus Multi-Slave Sim 10502" python modbus_sim.py
echo 启动完成，端口: 10502 (从机 1, 2, 3)
timeout /t 2 >nul
