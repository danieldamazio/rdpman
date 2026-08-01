package main

import (
	"embed"
	"log"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend
var assets embed.FS

// checkSingleInstance trava a execução se já houver um processo rodando
func checkSingleInstance() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex := kernel32.NewProc("CreateMutexW")
	
	name, _ := syscall.UTF16PtrFromString("rdpman_instance_lock")
	ret, _, _ := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(name)))
	
	// ERROR_ALREADY_EXISTS = 183
	if ret != 0 && syscall.GetLastError() == syscall.Errno(183) {
		log.Fatal("RDPMan já está em execução.")
	}
}

func main() {
	checkSingleInstance()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "rdpman",
		Width:  400,
		Height: 650,
		DisableResize: true,     // Impede redimensionamento lateral para manter o visual widget
		Frameless:     true,     // Janela sem bordas nativas
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 28, G: 28, B: 30, A: 255}, // #1C1C1E
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}