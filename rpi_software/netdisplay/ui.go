package netdisplay

import (
	"encoding/json"
	"log"
	"os"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"go.uber.org/zap"
)

type UI struct {
	asel *astilectron.Astilectron
	awin *astilectron.Window
	log  *zap.SugaredLogger
}

func NewUI(logger *zap.SugaredLogger) *UI {
	//Initialize astilectron
	var asel, err = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
		AppName:           "MKTNE NetDisplay",
		BaseDirectoryPath: "webcode",
	})
	log := logger.Named("UI")
	if err != nil {
		log.Fatal("Error initializing astilectron", zap.Error(err))
	}
	// Handle signals
	asel.HandleSignals()
	return &UI{asel: asel, log: log}
}
func (sui *UI) StartUI() {
	var err error
	// Start astilectron
	if err = sui.asel.Start(); err != nil {
		sui.log.Fatal("main: starting astilectron failed: %w", err)
	}

	if sui.awin, err = sui.asel.NewWindow("webcode/html/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(720),
		Width:  astikit.IntPtr(1080),
	}); err != nil {
		sui.log.Fatal("main: new window failed: %w", err)
	}
	// Create windows
	if err = sui.awin.Create(); err != nil {
		sui.log.Fatal("main: creating window failed: %w", err)
	}
}

func (sui *UI) Close() {
	sui.asel.Close()
}

func (sui *UI) OpenDevTools() {
	sui.awin.OpenDevTools()
}

func (sui *UI) Wait() {
	sui.asel.Wait()
}

func (sui *UI) UpdateUI(time string, activescreen string, numstrike uint8) error {
	type msg struct {
		Time   string `json:"time"`
		Screen string `json:"screen"`
		Strike uint8  `json:"strike"`
	}
	message, err := json.Marshal(msg{time, activescreen, numstrike})
	sui.awin.SendMessage(message)

	return err
}
