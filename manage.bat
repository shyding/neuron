@echo off
chcp 65001 >nul
title Neuron & CCTL Multi-Environment Manager
cd /d "%~dp0"
python manage.py
pause
