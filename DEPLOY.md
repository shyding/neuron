# Neuron 工业边缘网关部署指南（含全套 64 款工业协议驱动）

本分支集成了 Neuron IIoT 边缘网关及其全套 64 款工业驱动插件，可在任何 Linux / WSL2 x86_64 环境中开箱即用。

---

## 包含协议清单

### 1. PLC 与工业控制器
- **西门子（Siemens）**：S7 ISOTCP、S7-300/400、S5 Fetch-Write、PPI/MPI
- **三菱（Mitsubishi）**：MELSEC 1E、3E、4E、FX 系列
- **欧姆龙（Omron）**：FINS TCP、FINS UDP、Hostlink C-Mode、EtherNet/IP CIP
- **罗克韦尔（Rockwell / Allen-Bradley）**：ControlLogix / CompactLogix (5000)、MicroLogix (1400)、5500、EtherNet/IP、DF1
- **国内及其他品牌 PLC**：汇川（Inovance）、和利时（HollySys）、台达（Delta）、信捷（Xinje）、永宏（Fatek）、横河（Yokogawa）、松下（Panasonic）、富士（Fuji）、安川（Yaskawa）等

### 2. 数控机床与机器人（CNC & Robotics）
- **数控机床**：发那科（FANUC FOCAS）、海德汉（HEIDENHAIN CNC）、凯恩帝（KND CNC）、Brother CNC、Mazak CNC
- **工业机器人**：KUKA 工业机器人驱动

### 3. 标准工业与行业协议
- **通用工业**：OPC UA Client、OPC UA Server、Beckhoff TwinCAT ADS、CODESYS V3、Modbus（TCP / RTU / ASCII / QH）
- **电力能源与楼宇**：IEC 60870-5-104、DNP 3.0、DLT645-1997 / 2007（电表）、BACnet/IP、KNXnet/IP
- **半导体**：SECS/GEM（HSMS 通信标准）

### 4. 北向物联网与流处理
- MQTT、MQTT v2、Kafka、AWS IoT、Azure IoT、Sparkplug B、WebSocket、eKuiper 流式分析引擎

---

## 一键快速部署

### 方式 1：控制台直接启动
进入 `bin/` 目录运行启动脚本：
```bash
cd bin
chmod +x start.sh neuron
./start.sh
```

### 方式 2：注册为系统服务（推荐）
在 `bin/` 目录下以 root / sudo 权限运行服务安装脚本：
```bash
cd bin
sudo chmod +x install_service.sh start.sh neuron
sudo ./install_service.sh
```
安装后支持标准的 systemd 服务控制：
- 查看状态：`sudo systemctl status neuron`
- 重启服务：`sudo systemctl restart neuron`
- 停止服务：`sudo systemctl stop neuron`

---

## Web 访问说明

服务启动后，通过浏览器访问：
- **Web 控制台地址**：`http://localhost:7000`
- **默认登录用户名**：`admin`
- **默认登录密码**：`0000`

进入控制台后：
- 在 **配置 -> 插件** 中可查看到所有已加载的工业协议驱动。
- 在 **配置 -> 南向设备** 中点击 **添加设备** 即可自由选择所有 55 款南向驱动并进行测点配置。
