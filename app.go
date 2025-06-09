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
