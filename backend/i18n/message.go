package i18n

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// MessageKey 用于标识消息类型
type MessageKey string

const (
	// 错误消息
	ErrOpenRegistry       MessageKey = "err_open_registry"
	ErrReadProxyEnable   MessageKey = "err_read_proxy_enable"
	ErrReadProxyServer   MessageKey = "err_read_proxy_server"
	ErrWriteRegistry     MessageKey = "err_write_registry"
	ErrExecutePowerShell MessageKey = "err_execute_powershell"
	ErrStopDNS          MessageKey = "err_stop_dns"
	ErrStartDNS         MessageKey = "err_start_dns"
	ErrResetTCPIP       MessageKey = "err_reset_tcpip"
	ErrResetWinsock     MessageKey = "err_reset_winsock"
	ErrResetProxy       MessageKey = "err_reset_proxy"
	ErrFlushDNS         MessageKey = "err_flush_dns"
	ErrPingTest         MessageKey = "err_ping_test"
	ErrPingFailed       MessageKey = "err_ping_failed"
	ErrGetCurrentIP     MessageKey = "err_get_current_ip"
	ErrBasicNetworkFix  MessageKey = "err_basic_network_fix"
	ErrReleaseIP        MessageKey = "err_release_ip"
	ErrRenewIP          MessageKey = "err_renew_ip"
	ErrResetFirewall    MessageKey = "err_reset_firewall"
	ErrGeneric          MessageKey = "err_generic"

	// 成功消息
	SuccessDisableProxy    MessageKey = "success_disable_proxy"
	SuccessDisableProxyPS  MessageKey = "success_disable_proxy_ps"
	SuccessResetProxy      MessageKey = "success_reset_proxy"
	SuccessFlushDNS        MessageKey = "success_flush_dns"
	SuccessResetTCPIP      MessageKey = "success_reset_tcpip"
	SuccessResetWinsock    MessageKey = "success_reset_winsock"
	SuccessRestartDNS      MessageKey = "success_restart_dns"
	SuccessPingTest        MessageKey = "success_ping_test"
	SuccessPingTestWithDelay MessageKey = "success_ping_test_with_delay"
	SuccessGetCurrentIP    MessageKey = "success_get_current_ip"
	SuccessBasicNetworkFix MessageKey = "success_basic_network_fix"
	SuccessReleaseRenewIP  MessageKey = "success_release_renew_ip"
	SuccessResetFirewall   MessageKey = "success_reset_firewall"
)

var messages = map[string]map[MessageKey]string{
	"zh": {
		ErrOpenRegistry:       "无法打开注册表键: %v",
		ErrReadProxyEnable:    "读取代理启用状态失败: %v",
		ErrReadProxyServer:    "读取代理服务器地址失败: %v",
		ErrWriteRegistry:      "写入注册表值失败: %v",
		ErrExecutePowerShell:  "PowerShell执行出错: %v\n输出: %s",
		ErrStopDNS:           "停止DNS服务时出错: %v\n输出: %s",
		ErrStartDNS:          "启动DNS服务时出错: %v\n输出: %s",
		ErrResetTCPIP:        "重置IP时出错: %v\n输出: %s",
		ErrResetWinsock:      "重置Winsock时出错: %v\n输出: %s",
		ErrResetProxy:        "重置系统代理时出错: %v\n输出: %s",
		ErrFlushDNS:          "清除DNS缓存时出错: %v\n输出: %s",
		ErrPingTest:          "Ping测试执行失败 %s: %v",
		ErrPingFailed:        "Ping测试失败: 无法连接到 %s",
		ErrGetCurrentIP:      "获取当前IP地址失败: %v",
		ErrBasicNetworkFix:   "基础网络修复时出错: %v",
		ErrReleaseIP:         "释放IP地址时出错: %v\n输出: %s",
		ErrRenewIP:           "重新获取IP地址时出错: %v\n输出: %s",
		ErrResetFirewall:     "重置防火墙时出错: %v\n输出: %s",
		ErrGeneric:           "失败: %v",

		SuccessDisableProxy:   "成功: 已通过直接修改注册表关闭系统代理",
		SuccessDisableProxyPS: "成功: 已通过PowerShell命令关闭系统代理",
		SuccessResetProxy:     "成功: 系统代理已重置",
		SuccessFlushDNS:       "成功: DNS缓存已清除",
		SuccessResetTCPIP:     "成功: TCP/IP栈已重置（可能需要重启）",
		SuccessResetWinsock:   "成功: Winsock协议已重置（可能需要重启）",
		SuccessRestartDNS:     "成功: DNS客户端缓存服务已重启",
		SuccessPingTest:       "成功: 网络连通性正常，可以访问 %s",
		SuccessPingTestWithDelay: "成功: 网络连通性正常，可以访问 %s，延迟: %s",
		SuccessGetCurrentIP:   "当前IP地址: %s",
		SuccessReleaseRenewIP: "成功: IP地址已释放并重新获取",
		SuccessResetFirewall:  "成功: 防火墙设置已重置为默认状态",
	},
	"en": {
		ErrOpenRegistry:       "Failed to open registry key: %v",
		ErrReadProxyEnable:    "Failed to read proxy enable status: %v",
		ErrReadProxyServer:    "Failed to read proxy server address: %v",
		ErrWriteRegistry:      "Failed to write registry value: %v",
		ErrExecutePowerShell:  "PowerShell execution error: %v\nOutput: %s",
		ErrStopDNS:           "Error stopping DNS service: %v\nOutput: %s",
		ErrStartDNS:          "Error starting DNS service: %v\nOutput: %s",
		ErrResetTCPIP:        "Error resetting IP: %v\nOutput: %s",
		ErrResetWinsock:      "Error resetting Winsock: %v\nOutput: %s",
		ErrResetProxy:        "Error resetting system proxy: %v\nOutput: %s",
		ErrFlushDNS:          "Error flushing DNS cache: %v\nOutput: %s",
		ErrPingTest:          "Ping test execution failed %s: %v",
		ErrPingFailed:        "Ping test failed: Cannot connect to %s",
		ErrGetCurrentIP:      "Failed to get current IP address: %v",
		ErrBasicNetworkFix:   "Error during basic network fix: %v",
		ErrReleaseIP:         "Error releasing IP address: %v\nOutput: %s",
		ErrRenewIP:           "Error renewing IP address: %v\nOutput: %s",
		ErrResetFirewall:     "Error resetting firewall: %v\nOutput: %s",
		ErrGeneric:           "Failed: %v",

		SuccessDisableProxy:   "Success: System proxy disabled via direct registry modification",
		SuccessDisableProxyPS: "Success: System proxy disabled via PowerShell command",
		SuccessResetProxy:     "Success: System proxy reset",
		SuccessFlushDNS:       "Success: DNS cache cleared",
		SuccessResetTCPIP:     "Success: TCP/IP stack reset (may require restart)",
		SuccessResetWinsock:   "Success: Winsock protocol reset (may require restart)",
		SuccessRestartDNS:     "Success: DNS client cache service restarted",
		SuccessPingTest:       "Success: Network connectivity is normal, can access %s",
		SuccessPingTestWithDelay: "Success: Network connectivity is normal, can access %s, delay: %s",
		SuccessGetCurrentIP:   "Current IP address: %s",
		SuccessReleaseRenewIP: "Success: IP address released and renewed",
		SuccessResetFirewall:  "Success: Firewall settings reset to default",
	},
}

var currentLocale = getSystemLocale()

// GetCurrentLocale 获取当前语言设置
func GetCurrentLocale() string {
	return currentLocale
}

// SetLocale 设置当前语言
func SetLocale(locale string) string {
	if _, ok := messages[locale]; ok {
		currentLocale = locale
		return currentLocale
	}
	return "zh" // 默认返回中文
}

// getSystemLocale 获取系统语言设置
func getSystemLocale() string {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\International`, registry.QUERY_VALUE)
	if err != nil {
		return "zh" // 默认中文
	}
	defer key.Close()

	locale, _, err := key.GetStringValue("LocaleName")
	if err != nil {
		return "zh"
	}

	if locale == "zh-CN" || locale == "zh-TW" || locale == "zh-HK" {
		return "zh"
	}
	return "en"
}

// GetMessage 获取指定key的消息
func GetMessage(key MessageKey, args ...interface{}) string {
	if msg, ok := messages[currentLocale][key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}
	return fmt.Sprintf("Unknown message key: %s", key)
}