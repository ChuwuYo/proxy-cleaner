# ProxyCleaner <img src="frontend/public/appicon.png" alt="ProxyCleaner Logo" width="40" height="40" style="vertical-align: middle;">

[中文](README.md) | [English](README_EN.md)

## 项目简介

ProxyCleaner 是一个基于 Wails 框架开发的桌面应用程序，旨在帮助用户快速清理和管理 Windows 系统的网络配置。

## 主要功能

- **代理管理**：
  - 查看当前系统代理状态和配置
  - 快速禁用系统代理（支持直接修改注册表 和 PowerShell修改注册表 两种方式）
  - 重置系统代理设置

- **网络连通性测试**：
  - 支持自定义主机地址的 Ping 测试
  - 显示网络延迟信息，包括每次响应的延迟时间

- **网络修复工具**：
  - DNS 相关：
    - 清除 DNS 缓存
    - 重启 DNS 客户端缓存服务
  - 网络协议重置：
    - 重置 TCP/IP 栈
    - 重置 Winsock 协议

## 使用警告

### 权限说明
本程序的以下功能在执行时会自动请求管理员权限：
- 重置 TCP/IP 栈
- 重置 Winsock 协议
- 重启 DNS 客户端缓存服务

当执行这些操作时，系统会弹出UAC权限提升请求窗口，请点击"是"允许操作执行。

### 危险操作警告
以下操作可能会暂时中断网络连接，并且需要重启计算机才能完全生效：

1. **重置 TCP/IP 栈**
   - 此操作会重置所有网络适配器的 TCP/IP 设置
   - 可能会导致网络连接暂时中断
   - 建议在执行此操作前保存所有工作
   - 操作后可能需要重启计算机

2. **重置 Winsock 协议**
   - 此操作会重置 Windows 网络编程接口
   - 可能会影响所有网络相关程序
   - 建议在执行此操作前关闭所有网络应用
   - 操作后可能需要重启计算机

3. **重启 DNS 客户端缓存服务**
   - 此操作会暂时中断 DNS 解析服务
   - 可能会导致域名解析暂时失败
   - 建议在浏览器关闭的情况下执行

## 开发指南

### 环境准备

确保您已安装以下工具：

- Go (1.18 或更高版本)
- Node.js (推荐 LTS 版本)
- Wails CLI (通过 `go install github.com/wailsapp/wails/v2/cmd/wails@latest` 安装)

### 运行项目

在 `proxy-cleaner\frontend` 目录下运行以下命令安装前端依赖：

```bash
npm install
```

在项目根目录下运行以下命令以启动开发模式：

```bash
wails dev
```

这将启动一个 Vite 开发服务器，并提供前端代码的快速热重载。

## 构建项目

### 使用构建脚本（推荐）

项目提供了自动化构建脚本，可以自动更新版本信息并构建应用：

```powershell
.\build.ps1 -Version "0.0.8"
```

这个脚本会：
- 自动更新 `wails.json` 中的版本信息
- 执行构建命令并传递正确的版本号
- 确保 Windows 属性和应用标题显示相同的版本号

### 手动构建

如果需要手动构建，请先确保 `wails.json` 中的版本号正确，然后运行：

```bash
wails build -ldflags="-X main.Version=0.0.8"
```

## 贡献

欢迎任何形式的贡献！如果您有任何建议或发现 Bug，请随时提交 Issue 或 Pull Request。
