@echo off
echo ===================================================
echo   Starting NeuronEX Industrial Edge Platform...
echo ===================================================
cd /d "%~dp0"
start "NeuronEX" bin\neuronex.exe alone
echo [OK] NeuronEX service started!
echo Web Dashboard available at: http://localhost:8086/web/common/
echo Default credentials: admin / 0000
pause
