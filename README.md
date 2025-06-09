# ProxyCleaner

## 项目简介

ProxyCleaner 是一个基于 Wails 框架开发的桌面应用程序，旨在帮助用户清理代理残留设置。

![image](https://github.com/user-attachments/assets/831e722d-0e24-41c7-bcd9-b63c8791dbd3)


## 主要功能

- **代理清理**：方便地禁用系统代理,清理系统中可能存在的无效或残留代理配置。

## 开发指南

### 环境准备

确保您已安装以下工具：

- Go (1.18 或更高版本)
- Node.js (推荐 LTS 版本)
- Wails CLI (本项目暂为V2版本)

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
