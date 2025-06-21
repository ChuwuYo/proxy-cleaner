# ProxyCleaner <img src="frontend/public/appicon.png" alt="ProxyCleaner Logo" width="40" height="40" style="vertical-align: middle;">

[中文](README.md) | [English](README_EN.md)

## Project Introduction

ProxyCleaner is a desktop application developed based on the Wails framework, designed to help users quickly clean up and manage Windows system network configurations.

## Main Features

- **Proxy Management**:
  - View current system proxy status and configuration
  - Quickly disable system proxy (supports both direct registry modification and PowerShell registry modification)
  - Reset system proxy settings

- **Network Connectivity Test**:
  - Ping test with custom host address support
  - Display network latency information, including response time for each ping

- **Network Repair Tools**:
  - DNS related:
    - Flush DNS cache
    - Restart DNS client cache service
  - Network protocol reset:
    - Reset TCP/IP stack
    - Reset Winsock protocol

## Usage Warnings

### Permission Notes
The following features of this program will automatically request administrator privileges when executed:
- Reset TCP/IP stack
- Reset Winsock protocol
- Restart DNS client cache service

When performing these operations, the system will display a UAC privilege elevation request window. Please click "Yes" to allow the operation to proceed.

### Dangerous Operation Warnings
The following operations may temporarily interrupt network connections and require a computer restart to fully take effect:

1. **Reset TCP/IP Stack**
   - This operation resets TCP/IP settings for all network adapters
   - May cause temporary network connection interruption
   - It is recommended to save all work before performing this operation
   - A computer restart may be required after the operation

2. **Reset Winsock Protocol**
   - This operation resets the Windows network programming interface
   - May affect all network-related programs
   - It is recommended to close all network applications before performing this operation
   - A computer restart may be required after the operation

3. **Restart DNS Client Cache Service**
   - This operation will temporarily interrupt DNS resolution services
   - May cause temporary domain name resolution failures
   - It is recommended to perform this with browsers closed

## Development Guide

### Environment Preparation

Ensure you have the following tools installed:

- Go (1.18 or higher)
- Node.js (LTS version recommended)
- Wails CLI (install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Running the Project

Run the following command in the `proxy-cleaner\frontend` directory to install frontend dependencies:

```bash
npm install
```

Run the following command in the project root directory to start development mode:

```bash
wails dev
```

This will start a Vite development server and provide fast hot reloading of frontend code.

## Building the Project

### Using Build Script (Recommended)

The project provides an automated build script that automatically updates version information and builds the application:

```powershell
.\build.ps1 -Version "0.0.8"
```

This script will:
- Automatically update version information in `wails.json`
- Execute the build command with the correct version number
- Ensure Windows properties and application title display the same version number

### Manual Build

If you need to build manually, first ensure the version number in `wails.json` is correct, then run:

```bash
wails build -ldflags="-X main.Version=0.0.8"
```

## Contribution

Contributions of any form are welcome! If you have any suggestions or find bugs, please feel free to submit an Issue or Pull Request.