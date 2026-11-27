@echo off
echo Stopping NeuronEX...
taskkill /F /IM neuronex.exe >nul 2>&1
echo [OK] Stopped.
pause
