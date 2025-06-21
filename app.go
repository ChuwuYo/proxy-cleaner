package main

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"proxy-cleaner/backend/i18n"
)

// 检查当前进程是否具有管理员权限
func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

// 检查命令执行是否需要管理员权限
func requiresAdmin(command string) bool {
	adminCommands := []string{
		"netsh", "net stop", "net start",
	}
	for _, cmd := range adminCommands {
		if strings.Contains(command, cmd) {
			return true
		}
	}
	return false
}

// App struct holds application context
type App struct {
	ctx context.Context
}

// GetCurrentLocale 获取当前语言设置
func (a *App) GetCurrentLocale() string {
	return i18n.GetCurrentLocale()
}

// SetLocale 设置当前语言
func (a *App) SetLocale(locale string) string {
	return i18n.SetLocale(locale)
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// executeAsAdmin runs a command with elevated privileges
func executeAsAdmin(command string, args ...string) ([]byte, error) {
	psCmd := fmt.Sprintf("Start-Process -FilePath '%s' -ArgumentList '%s' -Verb RunAs -Wait -PassThru",
		command, strings.Join(args, " "))
	cmd := exec.Command("powershell", "-Command", psCmd)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	return cmd.CombinedOutput()
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
		return ProxyStatus{Error: i18n.GetMessage(i18n.ErrOpenRegistry, err.Error())}
	}
	defer key.Close()

	// Read proxy enable status
	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		// If the value doesn't exist or there's an error, assume proxy is disabled.
		proxyEnable = 0
	}

	// Read proxy server address with proper error handling
	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil && err != registry.ErrNotExist {
		return ProxyStatus{Error: i18n.GetMessage(i18n.ErrReadProxyServer, err.Error())}
	}

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
		return fmt.Errorf(i18n.GetMessage(i18n.ErrOpenRegistry, err))
	}
	defer key.Close()

	var dwordValue uint32 = 0
	if enabled {
		dwordValue = 1
	}

	err = key.SetDWordValue("ProxyEnable", dwordValue)
	if err != nil {
		return fmt.Errorf(i18n.GetMessage(i18n.ErrWriteRegistry, err))
	}
	return nil
}

// DisableProxyDirectly disables the system proxy by directly modifying the registry.
func (a *App) DisableProxyDirectly() string {
	err := setProxyState(false)
	if err != nil {
		return i18n.GetMessage(i18n.ErrGeneric, err)
	}
	return i18n.GetMessage(i18n.SuccessDisableProxy)
}

// DisableProxyViaPowerShell disables the system proxy using a PowerShell command.
func (a *App) DisableProxyViaPowerShell() string {
	cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings' -Name ProxyEnable -Value 0")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Combine error message with PowerShell's output for better debugging.
		return i18n.GetMessage(i18n.ErrExecutePowerShell, err, string(output))
	}
	return i18n.GetMessage(i18n.SuccessDisableProxyPS)
}

// ResetSystemProxy 重置系统代理设置
func (a *App) ResetSystemProxy() string {
	cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return i18n.GetMessage(i18n.ErrGeneric, err, string(output))
	}
	return i18n.GetMessage(i18n.SuccessResetProxy)
}

// FlushDNSCache 清除 DNS 缓存
func (a *App) FlushDNSCache() string {
	cmd := exec.Command("ipconfig", "/flushdns")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return i18n.GetMessage(i18n.ErrGeneric, err, string(output))
	}
	return i18n.GetMessage(i18n.SuccessFlushDNS)
}

// ResetTCPIP 重置 TCP/IP 栈
func (a *App) ResetTCPIP() string {
	// 使用提升的权限执行命令
	output, err := executeAsAdmin("netsh", "int", "ip", "reset")
	if err != nil {
		return i18n.GetMessage(i18n.ErrResetTCPIP, err, string(output))
	}
	return i18n.GetMessage(i18n.SuccessResetTCPIP)
}

// ResetWinsock 重置 Winsock 协议
func (a *App) ResetWinsock() string {
	// 使用提升的权限执行命令
	output, err := executeAsAdmin("netsh", "winsock", "reset")
	if err != nil {
		return i18n.GetMessage(i18n.ErrResetWinsock, err, string(output))
	}
	return i18n.GetMessage(i18n.SuccessResetWinsock)
}

// RestartDNSService 重启 DNS 客户端缓存服务
func (a *App) RestartDNSService() string {
	// 先检查服务是否在运行
	checkCmd := exec.Command("sc", "query", "dnscache")
	checkCmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	checkOutput, _ := checkCmd.CombinedOutput()
	isRunning := strings.Contains(string(checkOutput), "RUNNING")

	// 停止服务
	out1, err1 := executeAsAdmin("net", "stop", "dnscache")
	if err1 != nil {
		if isRunning {
			return i18n.GetMessage(i18n.ErrStopDNS, err1, string(out1))
		}
		// 如果服务本来就没在运行，则继续尝试启动
	}

	// 启动服务
	out2, err2 := executeAsAdmin("net", "start", "dnscache")
	if err2 != nil {
		// 如果启动失败且服务之前在运行，尝试恢复原状态
		if isRunning {
			executeAsAdmin("net", "start", "dnscache") // 忽略恢复错误，因为已经处于错误状态
		}
		return i18n.GetMessage(i18n.ErrStartDNS, err2, string(out2))
	}
	
	return i18n.GetMessage(i18n.SuccessRestartDNS)
}

// PingTest 执行ping连通性测试
func (a *App) PingTest(host string) string {
	if host == "" {
		host = "bing.com"
	}
	
	cmd := exec.Command("ping", "-n", "4", host)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return i18n.GetMessage(i18n.ErrPingTest, host, err.Error())
	}
	
	// 解析ping结果和延迟信息
	result := string(output)
	if strings.Contains(result, "TTL=") {
		// 只匹配包含TTL的响应行中的延迟
		lines := strings.Split(result, "\n")
		var delays []string
		
		for _, line := range lines {
			if strings.Contains(line, "TTL=") {
				// 匹配该行中的延迟，支持 XXXms 和 <XXXms 格式
				reTime := regexp.MustCompile(`<?([0-9]+)ms`)
				matches := reTime.FindStringSubmatch(line)
				if len(matches) > 1 {
					if strings.Contains(line, "<"+matches[1]+"ms") {
						delays = append(delays, "<"+matches[1]+"ms")
					} else {
						delays = append(delays, matches[1]+"ms")
					}
				}
			}
		}
		
		if len(delays) > 0 {
			avgDelay := strings.Join(delays, ", ")
			return i18n.GetMessage(i18n.SuccessPingTestWithDelay, host, avgDelay)
		}
		return i18n.GetMessage(i18n.SuccessPingTest, host)
	} else {
		return i18n.GetMessage(i18n.ErrPingFailed, host)
	}
}