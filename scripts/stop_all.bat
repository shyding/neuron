@echo off
chcp 65001 >nul
echo ======================================================
echo          一键停止所有服务 (10502 / 7000 / 8085 / 8088)
echo ======================================================
echo 正在停止 Windows 8088 原生服务...
taskkill /F /IM neuronex.exe 2>nul

echo 正在停止 WSL 7000 / 8085 服务...
wsl -d Ubuntu -e bash -c "killall -9 neuron neuronex 2>/dev/null || true"

echo 正在停止 10502 Modbus 仿真器...
powershell -Command "Get-CimInstance Win32_Process | Where-Object { $_.CommandLine -like '*modbus_sim.py*' } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force }" 2>nul

echo ======================================================
echo 所有环境服务已停止。
echo ======================================================
pause
