package main

import (
	"embed"
	"fmt"
	"log" // 引入 log 包，用于在启动时打印信息

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// 这些变量在构建时通过 ldflags 被 Wails 自动注入,不需要赋值
var (
	Version   string
	Commit    string
	BuildTime string
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed all:frontend/public/appicon.png
var icon []byte

// BuildInfo 结构体用于一次性将所有构建信息传递给前端
type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

func main() {
	// 创建 app 实例
	app := NewApp()

	// 在程序启动时，将构建信息打印到控制台，方便调试和验证
	log.Printf("Proxy-Cleaner Version: %s, Commit: %s, Build Time: %s\n", Version, Commit, BuildTime)

	// 创建应用
	err := wails.Run(&options.App{
		// 程序标题动态显示版本号
		Title:  fmt.Sprintf("Proxy-Cleaner v%s", Version),
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WindowIsTranslucent: false,
			DisableWindowIcon:   false,
		},
		MinWidth:  600,
		MinHeight: 500,
		MaxWidth:  2560,
		MaxHeight: 1600,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// 前端调用,获取所有构建信息
func (a *App) GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}
}
