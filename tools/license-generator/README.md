# NeuronEX 许可证与激活码生成配套工具 (License & Activation Code Generator)

本工具为自主实现的 NeuronEX 工业边缘计算网关专属配套工具，支持独立生成官方商业/企业版许可证文件（`.lic`）和一键激活码（Activation Code）。

---

## 功能特性

1. **生成激活码 (Activation Code)**：
   - 格式标准：`NX-<Base32Payload>-<HMACChecksum>`
   - 可直接在 NeuronEX Web 界面 `系统管理 -> 许可证 -> 方式 2 (激活码)` 中输入并一键激活。
2. **生成许可证文件 (`.lic`)**：
   - 生成符合 NeuronEX 企业规范的完整授权描述与数字防伪签名。
   - 可在 Web 界面 `方式 1 (上传许可证)` 中直接导入生效。
3. **消除所有试用限制 ("不能自己限制自己")**：
   - 生成的许可证默认永久有效（2999年到期）。
   - 核心采集点位支持无上限配置（默认 1,000,000 点）。
   - 默认开启全部 59 种工业协议插件（Modbus、OPC UA/DA、Siemens S7、Mitsubishi、Omron、BACnet 等）。
   - 默认完整包含 eKuiper 工业流式数据分析引擎，彻底消除 30 分钟试用中断限制。
4. **校验模式**：
   - 支持验证激活码有效性 (`--verify-code`) 和许可证文件完整性 (`--verify-file`)。

---

## 编译方法

在项目根目录下执行：

```bash
# 编译生成可执行二进制文件
go build -o bin/neuronex-licgen ./tools/license-generator
```

---

## 使用说明

### 1. 默认快速生成（生成企业版许可证及激活码）
```bash
./bin/neuronex-licgen
```
输出示例：
```
================================================================================
             NeuronEX License & Activation Code Generator Tool
================================================================================
[Configuration]
  • Customer Name   : EMQ Technologies Inc.
  • Customer Contact: support@emqx.io
  • Tag Capacity    : 1000000 points
  • Expire Date     : 2999-01-01
  • Hardware Token  : NEURONEX-CORE-ENTERPRISE
--------------------------------------------------------------------------------

🔑 [Generated Activation Code]:

   NX-IVWUK4RUMF2A4Z3JNZXXE2LOEB2W63TFON2W2Y2QOBZGS43VOB2HE4DP-C4927289

   --> Copy and paste this code into NeuronEX Web UI (方式 2: 激活码) to activate!

📄 [Generated License File]: neuronex.lic
   --> Upload this file into NeuronEX Web UI (方式 1: 上传许可证) to activate!
================================================================================
```

### 2. 为指定客户定制生成企业许可证
```bash
./bin/neuronex-licgen \
  -customer "上海某某科技有限公司" \
  -contact "tech@example.com" \
  -tags 500000 \
  -expire "2035-12-31" \
  -out custom.lic
```

### 3. 校验激活码
```bash
./bin/neuronex-licgen -verify-code NX-IVWUK4RUMF2A4Z3JNZXXE2LOEB2W63TFON2W2Y2QOBZGS43VOB2HE4DP-C4927289
```

### 4. 校验许可证文件
```bash
./bin/neuronex-licgen -verify-file custom.lic
```
