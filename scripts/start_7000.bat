@echo off
chcp 65001 >nul
echo [7000] 正在启动 WSL Neuron C Core (端口 7000)...
start "WSL Neuron C Core (7000)" wsl.exe -d Ubuntu bash -c "cd /opt/neuronex-4.0.2 && ./bin/neuron"
echo 启动完成，端口: 7000
timeout /t 2 >nul
