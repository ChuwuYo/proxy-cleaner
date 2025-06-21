param (
    [Parameter(Mandatory=$true)]
    [string]$Version
)

# 读取wails.json文件
$wailsJsonPath = "wails.json"
if (-not (Test-Path $wailsJsonPath)) {
    Write-Error "找不到 wails.json 文件"
    exit 1
}

# 读取并解析JSON
$wailsConfig = Get-Content $wailsJsonPath -Raw | ConvertFrom-Json

# 更新版本号
$wailsConfig.info.productVersion = $Version

# 写回JSON文件
$wailsConfig | ConvertTo-Json -Depth 10 | Out-File -FilePath $wailsJsonPath -Encoding utf8

Write-Host "已更新 wails.json 文件，版本号: $Version"

# 执行wails构建命令
$buildCommand = "wails build -ldflags=`"-X main.Version=$Version`""
Write-Host "执行构建命令: $buildCommand"
Invoke-Expression $buildCommand

Write-Host "构建完成！"