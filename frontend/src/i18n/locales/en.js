export default {
  title: "ProxyCleaner",
  subtitle: "A tool for quickly cleaning Windows system proxy residuals",
  status: {
    title: "Current ProxyEnable Status",
    error: "Failed to get ProxyEnable status: {msg}",
    enabled: "Enabled",
    disabled: "Disabled",
    address: "Address",
    notSet: "Not Set",
    notEffective: "Not Effective"
  },
  operations: {
    title: "ProxyEnable Operations",
    direct: "1. Direct Registry Modification",
    powershell: "2. PowerShell Registry Modification"
  },
  basicRepair: {
    title: "Basic Network Repair",
    resetProxy: "Reset System Proxy",
    flushDNS: "Flush DNS Cache"
  },
  advancedReset: {
    title: "Advanced Network Reset",
    resetTCPIP: "Reset TCP/IP Stack",
    resetWinsock: "Reset Winsock Protocol",
    restartDNS: "Restart DNS Service"
  },
  logs: {
    title: "Operation Logs",
    refreshing: "Getting proxy status...",
    updateSuccess: "Proxy status updated.",
    statusRefreshed: "Proxy status refreshed",
    backendError: "Backend call failed: {msg}",
    executing: "Executing operation...",
    directModifying: "Attempting to modify registry directly...",
    psModifying: "Attempting to modify registry via PowerShell...",
    resettingProxy: "Resetting system proxy settings...",
    flushingDNS: "Flushing DNS cache...",
    resettingTCPIP: "Resetting TCP/IP stack...",
    resettingWinsock: "Resetting Winsock protocol...",
    restartingDNS: "Restarting DNS service..."
  },
  common: {
    success: "Success",
    failed: "Failed",
    error: "Error"
  }
}