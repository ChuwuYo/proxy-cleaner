export default {
  title: "ProxyCleaner",
  subtitle: "一个用于快速清理 Windows 系统代理残留的工具",
  status: {
    title: "当前ProxyEnable状态",
    error: "获取ProxyEnable状态失败: {msg}",
    enabled: "已启用",
    disabled: "已禁用",
    unknown: "未知",
    address: "地址",
    notSet: "未设置",
    notEffective: "未生效"
  },
  operations: {
    title: "ProxyEnable 禁用操作",
    direct: "1. 直接修改注册表",
    powershell: "2. PowerShell修改注册表"
  },
  basicRepair: {
    title: "网络基础修复",
    resetProxy: "重置系统代理设置",
    flushDNS: "清除 DNS 缓存",
    releaseRenewIP: "释放并重新获取IP"
  },
  advancedReset: {
    title: "高级网络重置",
    resetTCPIP: "重置 TCP/IP 栈",
    resetWinsock: "重置 Winsock 协议",
    restartDNS: "重启 DNS 服务",
    resetFirewall: "重置防火墙设置"
  },
  connectivity: {
    title: "网络连通性测试",
    currentIP: "当前IP地址",
    placeholder: "输入要测试的主机地址（默认：bing.com）",
    test: "开始测试",
    emptyHost: "请输入要测试的主机地址"
  },
  logs: {
    title: "运行日志",
    refreshing: "正在获取代理状态...",
    updateSuccess: "代理状态已更新。",
    statusRefreshed: "代理状态已刷新",
    backendError: "调用后端失败: {msg}",
    executing: "正在执行操作...",
    directModifying: "正在尝试直接修改注册表...",
    psModifying: "正在尝试通过PowerShell修改注册表...",
    resettingProxy: "正在重置系统代理设置...",
    flushingDNS: "正在清除 DNS 缓存...",
    resettingTCPIP: "正在重置 TCP/IP 栈...",
    resettingWinsock: "正在重置 Winsock 协议...",
    restartingDNS: "正在重启 DNS 服务...",
    resettingFirewall: "正在重置防火墙设置...",
    releaseRenewingIP: "正在释放并重新获取IP...",
    pingTesting: "正在测试连接到 {host}...",
    gettingIP: "正在获取当前IP地址...",
    ipUpdated: "IP地址已更新"
  },
  common: {
    success: "成功",
    failed: "失败",
    error: "错误"
  }
}