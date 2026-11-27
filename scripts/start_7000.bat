@echo off
chcp 65001 >nul
echo [7000] 正在启动 WSL Neuron C Core (端口 7000)...
wsl -d Ubuntu -e bash -c "cd /opt/neuronex-4.0.2 && nohup ./bin/neuron > /tmp/neuron.log 2>&1 &"
echo 启动完成，请检查端口 7000。
timeout /t 2 >nul
