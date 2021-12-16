package main

import (
	"log"
	"os"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/jax-b/ModulerKTNE/rpi_software/netdisplay"
	"go.uber.org/zap"
)

func main() {
	zaplogger := netdisplay.NewLogger()
	appconfig := netdisplay.NewConfig(zaplogger)
	appconfig.Load()
	zaplogger.Info("Config loaded")
	//Initialize astilectron
	var a, err = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:           "MKTNE NetDisplay",
		BaseDirectoryPath: "webcode",
	})
	if err != nil {
		zaplogger.Fatal("Error initializing astilectron", zap.Error(err))
	}
	defer a.Close()
	// Handle signals
	a.HandleSignals()

	// Start astilectron
	if err = a.Start(); err != nil {
		zaplogger.Fatal("main: starting astilectron failed: %w", err)
	}

	var w *astilectron.Window
	if w, err = a.NewWindow("webcode/html/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(720),
		Width:  astikit.IntPtr(1080),
	}); err != nil {
		zaplogger.Fatal("main: new window failed: %w", err)
	}
	// Create windows
	if err = w.Create(); err != nil {
		zaplogger.Fatal("main: creating window failed: %w", err)
	}
	w.OpenDevTools()
	// Blocking pattern
	a.Wait()
}
