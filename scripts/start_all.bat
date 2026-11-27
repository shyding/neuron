@echo off
chcp 65001 >nul
echo ======================================================
echo          一键启动所有环境服务 (10502 / 7000 / 8085 / 8088)
echo ======================================================
call "%~dp0start_modbus_sim.bat"
call "%~dp0start_7000.bat"
call "%~dp0start_8085.bat"
call "%~dp0start_8088.bat"
echo ======================================================
echo 所有启动指令已执行完毕！
echo Web 控制台:
echo   - 8088: http://localhost:8088/web/common/ (admin/0000)
echo   - 8085: http://localhost:8085/web/common/ (admin/0000)
echo   - 9100: http://localhost:9100/#/device/instance (admin/admin)
echo ======================================================
pause
