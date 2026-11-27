@echo off
chcp 65001 >nul
echo [8085] 正在启动 WSL 官方 NeuronEX (端口 8085)...
start "WSL Official NeuronEX (8085)" wsl.exe -d Ubuntu bash -c "cd /opt/neuronex-4.0.2 && ./bin/neuronex alone"
echo 启动完成，Web 界面: http://localhost:8085/web/common/
timeout /t 2 >nul
