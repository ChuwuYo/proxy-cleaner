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
	ErrReadProxyServer   MessageKey = "err_read_proxy_server"
	ErrWriteRegistry     MessageKey = "err_write_registry"
	ErrExecutePowerShell MessageKey = "err_execute_powershell"
	ErrStopDNS          MessageKey = "err_stop_dns"
	ErrStartDNS         MessageKey = "err_start_dns"
	ErrResetTCPIP       MessageKey = "err_reset_tcpip"
	ErrResetWinsock     MessageKey = "err_reset_winsock"
	ErrPingTest         MessageKey = "err_ping_test"
	ErrPingFailed       MessageKey = "err_ping_failed"
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
)

var messages = map[string]map[MessageKey]string{
	"zh": {
		ErrOpenRegistry:       "无法打开注册表键: %v",
		ErrReadProxyServer:    "读取代理服务器地址失败: %v",
		ErrWriteRegistry:      "写入注册表值失败: %v",
		ErrExecutePowerShell:  "PowerShell执行出错: %v\n输出: %s",
		ErrStopDNS:           "停止DNS服务时出错: %v\n输出: %s",
		ErrStartDNS:          "启动DNS服务时出错: %v\n输出: %s",
		ErrResetTCPIP:        "重置IP时出错: %v\n输出: %s",
		ErrResetWinsock:      "重置Winsock时出错: %v\n输出: %s",
		ErrPingTest:          "Ping测试执行失败 %s: %v",
		ErrPingFailed:        "Ping测试失败: 无法连接到 %s",
		ErrGeneric:           "失败: %v",

		SuccessDisableProxy:   "成功: 已通过直接修改注册表关闭系统代理。",
		SuccessDisableProxyPS: "成功: 已通过PowerShell命令关闭系统代理。",
		SuccessResetProxy:     "成功: 系统代理已重置。",
		SuccessFlushDNS:       "成功: DNS 缓存已清除。",
		SuccessResetTCPIP:     "成功: TCP/IP 栈已重置。警告：此操作可能需要重启计算机才能生效。",
		SuccessResetWinsock:   "成功: Winsock 协议已重置。警告：此操作可能需要重启计算机才能生效。",
		SuccessRestartDNS:     "成功: DNS 客户端缓存服务已重启。",
		SuccessPingTest:       "成功: 网络连通性正常，可以访问 %s",
		SuccessPingTestWithDelay: "成功: 网络连通性正常，可以访问 %s，延迟: %s",
	},
	"en": {
		ErrOpenRegistry:       "Failed to open registry key: %v",
		ErrReadProxyServer:    "Failed to read proxy server address: %v",
		ErrWriteRegistry:      "Failed to write registry value: %v",
		ErrExecutePowerShell:  "PowerShell execution error: %v\nOutput: %s",
		ErrStopDNS:           "Error stopping DNS service: %v\nOutput: %s",
		ErrStartDNS:          "Error starting DNS service: %v\nOutput: %s",
		ErrResetTCPIP:        "Error resetting IP: %v\nOutput: %s",
		ErrResetWinsock:      "Error resetting Winsock: %v\nOutput: %s",
		ErrPingTest:          "Ping test execution failed %s: %v",
		ErrPingFailed:        "Ping test failed: Cannot connect to %s",
		ErrGeneric:           "Failed: %v",

		SuccessDisableProxy:   "Success: System proxy disabled via direct registry modification.",
		SuccessDisableProxyPS: "Success: System proxy disabled via PowerShell command.",
		SuccessResetProxy:     "Success: System proxy reset.",
		SuccessFlushDNS:       "Success: DNS cache cleared.",
		SuccessResetTCPIP:     "Success: TCP/IP stack reset. Warning: This operation may require a computer restart.",
		SuccessResetWinsock:   "Success: Winsock protocol reset. Warning: This operation may require a computer restart.",
		SuccessRestartDNS:     "Success: DNS client cache service restarted.",
		SuccessPingTest:       "Success: Network connectivity is normal, can access %s",
		SuccessPingTestWithDelay: "Success: Network connectivity is normal, can access %s, delay: %s",
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