package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"time"

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

// NewApp 创建一个新的 App 应用结构体
func NewApp() *App {
	return &App{}
}

// startup 在应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// executeAsAdmin 以提升的权限运行命令
func executeAsAdmin(command string, args ...string) ([]byte, error) {
	psCmd := fmt.Sprintf("Start-Process -FilePath '%s' -ArgumentList '%s' -Verb RunAs -Wait -PassThru",
		command, strings.Join(args, " "))
	cmd := exec.Command("powershell", "-Command", psCmd)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	return cmd.CombinedOutput()
}

// ProxyStatus 定义返回给前端的代理状态的结构体
type ProxyStatus struct {
	Enabled bool   `json:"enabled"`
	Server  string `json:"server"`
	Error   string `json:"error"`
}

// OperationResult 定义操作结果的结构体，包含状态码和消息
type OperationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const regKeyPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

// GetProxyStatus 获取当前系统代理设置
func (a *App) GetProxyStatus() ProxyStatus {
	key, err := registry.OpenKey(registry.CURRENT_USER, regKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return ProxyStatus{Error: i18n.GetMessage(i18n.ErrOpenRegistry, err.Error())}
	}
	defer key.Close()

	// 读取代理启用状态，统一错误处理策略
	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil && err != registry.ErrNotExist {
		return ProxyStatus{Error: i18n.GetMessage(i18n.ErrReadProxyEnable, err.Error())}
	}
	if err == registry.ErrNotExist {
		proxyEnable = 0
	}

	// 读取代理服务器地址，统一错误处理策略
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

// setProxyState 是一个内部辅助函数，用于设置代理状态
func setProxyState(enabled bool) string {
	key, err := registry.OpenKey(registry.CURRENT_USER, regKeyPath, registry.SET_VALUE)
	if err != nil {
		return i18n.GetMessage(i18n.ErrOpenRegistry, err.Error())
	}
	defer key.Close()

	var dwordValue uint32 = 0
	if enabled {
		dwordValue = 1
	}

	err = key.SetDWordValue("ProxyEnable", dwordValue)
	if err != nil {
		return i18n.GetMessage(i18n.ErrWriteRegistry, err.Error())
	}
	return ""
}

// DisableProxyDirectly 通过直接修改注册表禁用系统代理
func (a *App) DisableProxyDirectly() OperationResult {
	errorMsg := setProxyState(false)
	if errorMsg != "" {
		return OperationResult{
			Success: false,
			Message: errorMsg,
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessDisableProxy),
	}
}

// DisableProxyViaPowerShell 使用 PowerShell 命令禁用系统代理
func (a *App) DisableProxyViaPowerShell() OperationResult {
	cmd := exec.Command("powershell", "-Command", "Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings' -Name ProxyEnable -Value 0")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 将错误消息与 PowerShell 的输出相结合，以便更好地进行调试。
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrExecutePowerShell, err.Error(), string(output)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessDisableProxyPS),
	}
}

// ResetSystemProxy 重置系统代理设置
func (a *App) ResetSystemProxy() OperationResult {
	cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrResetProxy, err.Error(), string(output)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessResetProxy),
	}
}

// FlushDNSCache 清除 DNS 缓存
func (a *App) FlushDNSCache() OperationResult {
	cmd := exec.Command("ipconfig", "/flushdns")
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrFlushDNS, err.Error(), string(output)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessFlushDNS),
	}
}

// ResetTCPIP 重置 TCP/IP 栈
func (a *App) ResetTCPIP() OperationResult {
	// 使用提升的权限执行命令
	output, err := executeAsAdmin("netsh", "int", "ip", "reset")
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrResetTCPIP, err.Error(), string(output)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessResetTCPIP),
	}
}

// ResetWinsock 重置 Winsock 协议
func (a *App) ResetWinsock() OperationResult {
	// 使用提升的权限执行命令
	output, err := executeAsAdmin("netsh", "winsock", "reset")
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrResetWinsock, err.Error(), string(output)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessResetWinsock),
	}
}

// RestartDNSService 重启 DNS 客户端缓存服务
// 优化错误处理逻辑，避免混淆的恢复操作
func (a *App) RestartDNSService() OperationResult {
	// 先检查服务是否在运行
	checkCmd := exec.Command("sc", "query", "dnscache")
	checkCmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	checkOutput, _ := checkCmd.CombinedOutput()
	isRunning := strings.Contains(string(checkOutput), "RUNNING")

	// 如果服务正在运行，先停止服务
	if isRunning {
		out1, err1 := executeAsAdmin("net", "stop", "dnscache")
		if err1 != nil {
			return OperationResult{
				Success: false,
				Message: i18n.GetMessage(i18n.ErrStopDNS, err1.Error(), string(out1)),
			}
		}
	}

	// 启动服务
	out2, err2 := executeAsAdmin("net", "start", "dnscache")
	if err2 != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrStartDNS, err2.Error(), string(out2)),
		}
	}
	
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessRestartDNS),
	}
}

// GetCurrentIP 获取当前设备的公网IP地址
func (a *App) GetCurrentIP() OperationResult {
	transport := &http.Transport{}
	
	// 读取系统代理设置
	proxyStatus := a.GetProxyStatus()
	if proxyStatus.Enabled && proxyStatus.Server != "" {
		proxyURL, err := url.Parse("http://" + proxyStatus.Server)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}
	
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: transport,
	}
	
	resp, err := client.Get("https://myip.addr.tools/")
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrGetCurrentIP, err.Error()),
		}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrGetCurrentIP, err.Error()),
		}
	}
	
	ip := strings.TrimSpace(string(body))
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessGetCurrentIP, ip),
	}
}

// PingTest 执行ping连通性测试
func (a *App) PingTest(host string) OperationResult {
	if host == "" {
		host = "bing.com"
	}
	
	cmd := exec.Command("ping", "-n", "4", host)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrPingTest, host, err.Error()),
		}
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
			return OperationResult{
				Success: true,
				Message: i18n.GetMessage(i18n.SuccessPingTestWithDelay, host, avgDelay),
			}
		}
		return OperationResult{
			Success: true,
			Message: i18n.GetMessage(i18n.SuccessPingTest, host),
		}
	} else {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrPingFailed, host),
		}
	}
}

// ResetFirewall 重置防火墙设置
func (a *App) ResetFirewall() OperationResult {
	output, err := executeAsAdmin("netsh", "advfirewall", "reset")
	if err != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrResetFirewall, err.Error(), string(output)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessResetFirewall),
	}
}

// ReleaseRenewIP 释放并重新获取IP地址
func (a *App) ReleaseRenewIP() OperationResult {
	// 先释放IP
	cmd1 := exec.Command("ipconfig", "/release")
	cmd1.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output1, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrReleaseIP, err1.Error(), string(output1)),
		}
	}

	// 再重新获取IP
	cmd2 := exec.Command("ipconfig", "/renew")
	cmd2.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	output2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		return OperationResult{
			Success: false,
			Message: i18n.GetMessage(i18n.ErrRenewIP, err2.Error(), string(output2)),
		}
	}
	return OperationResult{
		Success: true,
		Message: i18n.GetMessage(i18n.SuccessReleaseRenewIP),
	}
}
