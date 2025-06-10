package main

import (
	"context"
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

// App struct holds application context
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ProxyStatus defines the structure for returning proxy status to the frontend.
type ProxyStatus struct {
	Enabled bool   `json:"enabled"`
	Server  string `json:"server"`
	Error   string `json:"error"`
}

const regKeyPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

// GetProxyStatus retrieves the current system proxy settings.
func (a *App) GetProxyStatus() ProxyStatus {
	key, err := registry.OpenKey(registry.CURRENT_USER, regKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return ProxyStatus{Error: "无法打开注册表键: " + err.Error()}
	}
	defer key.Close()

	// Read proxy enable status
	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		// If the value doesn't exist or there's an error, assume proxy is disabled.
		proxyEnable = 0
	}

	// Read proxy server address. Ignore errors as it might not be set.
	proxyServer, _, _ := key.GetStringValue("ProxyServer")

	// 如果地址包含 http:// 或 https:// 前缀，则移除
	if len(proxyServer) > 7 {
		if proxyServer[:7] == "http://" {
			proxyServer = proxyServer[7:]
		} else if len(proxyServer) > 8 && proxyServer[:8] == "https://" {
			proxyServer = proxyServer[8:]
		}
	}

	return ProxyStatus{
		Enabled: proxyEnable == 1,
		Server:  proxyServer,
	}
}

// setProxyState is an internal helper function to set the proxy state.
func setProxyState(enabled bool) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regKeyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("无法打开注册表键进行写入: %v", err)
	}
	defer key.Close()

	var dwordValue uint32 = 0
	if enabled {
		dwordValue = 1
	}

	err = key.SetDWordValue("ProxyEnable", dwordValue)
	if err != nil {
		return fmt.Errorf("写入注册表值失败: %v", err)
	}
	return nil
}

// DisableProxyDirectly disables the system proxy by directly modifying the registry.
func (a *App) DisableProxyDirectly() string {
	err := setProxyState(false)
	if err != nil {
		return fmt.Sprintf("失败: %v", err)
	}
	return "成功: 已通过直接修改注册表关闭系统代理。"
}

// DisableProxyViaPowerShell disables the system proxy using a PowerShell command.
func (a *App) DisableProxyViaPowerShell() string {
	cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings' -Name ProxyEnable -Value 0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Combine error message with PowerShell's output for better debugging.
		return fmt.Sprintf("失败: PowerShell执行出错: %v\n输出: %s", err, string(output))
	}
	return "成功: 已通过PowerShell命令关闭系统代理。"
}

// ResetSystemProxy 重置系统代理设置
func (a *App) ResetSystemProxy() string {
	cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("失败: %v\n输出: %s", err, string(output))
	}
	return "成功: 系统代理已重置。"
}

// FlushDNSCache 清除 DNS 缓存
func (a *App) FlushDNSCache() string {
	cmd := exec.Command("ipconfig", "/flushdns")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("失败: %v\n输出: %s", err, string(output))
	}
	return "成功: DNS 缓存已清除。"
}

// ResetTCPIP 重置 TCP/IP 栈
func (a *App) ResetTCPIP() string {
	cmd := exec.Command("netsh", "int", "ip", "reset")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("失败: 重置IP时出错: %v\n输出: %s", err, string(output))
	}
	return "成功: TCP/IP 栈已重置。"
}

// ResetWinsock 重置 Winsock 协议
func (a *App) ResetWinsock() string {
	cmd := exec.Command("netsh", "winsock", "reset")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("失败: 重置Winsock时出错: %v\n输出: %s", err, string(output))
	}
	return "成功: Winsock 协议已重置。"
}

// RestartDNSService 重启 DNS 客户端缓存服务
func (a *App) RestartDNSService() string {
	cmd1 := exec.Command("net", "stop", "dnscache")
	out1, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		return fmt.Sprintf("失败: 停止DNS服务时出错: %v\n输出: %s", err1, string(out1))
	}
	
	cmd2 := exec.Command("net", "start", "dnscache")
	out2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		return fmt.Sprintf("失败: 启动DNS服务时出错: %v\n输出: %s", err2, string(out2))
	}
	return "成功: DNS 客户端缓存服务已重启。"
}
