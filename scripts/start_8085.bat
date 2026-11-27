@echo off
chcp 65001 >nul
echo [8085] 正在启动 WSL 官方 NeuronEX (端口 8085)...
wsl -d Ubuntu -e bash -c "cd /opt/neuronex-4.0.2 && nohup ./bin/neuronex alone > /tmp/neuronex.log 2>&1 &"
echo 启动完成，Web 界面: http://localhost:8085/web/common/
timeout /t 2 >nul
