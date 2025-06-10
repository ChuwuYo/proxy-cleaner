# ProxyCleaner

## 项目简介

ProxyCleaner 是一个基于 Wails 框架开发的桌面应用程序，旨在帮助用户快速清理和管理 Windows 系统的网络配置。

## 主要功能

- **代理管理**：
  - 查看当前系统代理状态和配置
  - 快速禁用系统代理（支持直接注册表和 PowerShell 两种方式）
  - 重置系统代理设置

- **网络修复工具**：
  - DNS 相关：
    - 清除 DNS 缓存
    - 重启 DNS 客户端缓存服务
  - 网络协议重置：
    - 重置 TCP/IP 栈
    - 重置 Winsock 协议

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

要构建可分发的生产模式软件包，请运行以下命令：

```bash
wails build
```

## 贡献

欢迎任何形式的贡献！如果您有任何建议或发现 Bug，请随时提交 Issue 或 Pull Request。
