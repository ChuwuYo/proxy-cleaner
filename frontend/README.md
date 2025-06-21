# ProxyCleaner 前端 <img src="public/appicon.png" alt="ProxyCleaner Logo" width="40" height="40" style="vertical-align: middle;">

## 项目简介

这是 ProxyCleaner 应用程序的前端部分，基于 Vue 3 和 Vite 构建，采用 Naive UI 作为组件库，提供直观的用户界面来管理系统网络配置。

## 功能界面

- **代理状态展示**：实时显示系统代理的启用状态和配置信息
- **快速操作区域**：提供代理禁用和网络修复的快捷操作按钮
- **网络连通性测试**：支持自定义主机的 Ping 测试及延迟显示
- **实时日志显示**：记录所有操作的执行状态和结果

## 技术栈

- **框架**: Vue 3
- **构建工具**: Vite
- **UI 组件库**: Naive UI
- **图标库**: `@vicons/ionicons5`

## 开发指南

### 依赖安装

在当前目录下运行以下命令安装前端依赖：

```bash
npm install
```

### 项目结构

- `src/App.vue`: 主应用组件，包含所有UI界面和业务逻辑
- `src/Root.vue`: 根组件，提供主题配置和全局消息提示
- `src/i18n/`: 国际化支持，包含中英文语言文件
- `src/assets/styles/`: 统一的样式管理系统
- `public/`: 静态资源目录，包含应用图标等文件
- `wailsjs/`: Wails 自动生成的 Go 后端接口绑定