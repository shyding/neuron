@echo off
chcp 65001 >nul
echo [8088] 正在启动 Windows 原生 Standalone NeuronEX (端口 8088)...
cd /d "%~dp0..\neuronex_windows_standalone"
start "NeuronEX 8088 Standalone" "%~dp0..\neuronex_windows_standalone\bin\neuronex.exe" alone
echo 启动完成，Web 界面: http://localhost:8088/web/common/
timeout /t 2 >nul
