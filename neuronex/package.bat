@echo off
cd /d "%~dp0"
python package.py %*
if errorlevel 1 pause
