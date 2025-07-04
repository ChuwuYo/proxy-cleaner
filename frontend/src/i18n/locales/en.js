export default {
  title: "ProxyCleaner",
  subtitle: "A tool for quickly cleaning Windows system proxy residuals",
  status: {
    title: "Current ProxyEnable Status",
    error: "Failed to get ProxyEnable status: {msg}",
    enabled: "Enabled",
    disabled: "Disabled",
    unknown: "Unknown",
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
    flushDNS: "Flush DNS Cache",
    releaseRenewIP: "Release & Renew IP"
  },
  advancedReset: {
    title: "Advanced Network Reset",
    resetTCPIP: "Reset TCP/IP Stack",
    resetWinsock: "Reset Winsock Protocol",
    restartDNS: "Restart DNS Service",
    resetFirewall: "Reset Firewall Settings"
  },
  connectivity: {
    title: "Network Connectivity Test",
    currentIP: "Current IP Address",
    placeholder: "Enter host address to test (default: bing.com)",
    test: "Start Test",
    emptyHost: "Please enter a host address to test"
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
    restartingDNS: "Restarting DNS service...",
    resettingFirewall: "Resetting firewall settings...",
    releaseRenewingIP: "Releasing and renewing IP...",
    pingTesting: "Testing connection to {host}...",
    gettingIP: "Getting current IP address...",
    ipUpdated: "IP address updated"
  },
  common: {
    success: "Success",
    failed: "Failed",
    error: "Error"
  }
}