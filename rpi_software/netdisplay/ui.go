package netdisplay

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"go.uber.org/zap"
)

type UI struct {
	Asel *astilectron.Astilectron
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
	asel.HandleSignals()
	return &UI{Asel: asel, log: log}
}
func (sui *UI) StartUI() {
	var err error
	// Start astilectron
	if err = sui.Asel.Start(); err != nil {
		sui.log.Fatal("main: starting astilectron failed: %w", err)
	}

	if sui.awin, err = sui.Asel.NewWindow("webcode/html/index.html", &astilectron.WindowOptions{
		Center:         astikit.BoolPtr(true),
		Height:         astikit.IntPtr(720),
		Width:          astikit.IntPtr(1080),
		Fullscreenable: astikit.BoolPtr(true),
	}); err != nil {
		sui.log.Fatal("main: new window failed: %w", err)
	}

	// Create windows
	if err = sui.awin.Create(); err != nil {
		sui.log.Fatal("main: creating window failed: %w", err)
	}
	m := sui.awin.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Menu"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Quit"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						sui.Asel.Stop()
						os.Exit(0)
						return
					},
				},
				{
					Label:   astikit.StrPtr("DevTools"),
					Checked: astikit.BoolPtr(false),
					Type:    astilectron.MenuItemTypeCheckbox,
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if *e.MenuItemOptions.Checked {
							sui.OpenDevTools()
						} else {
							sui.awin.CloseDevTools()
						}
						return
					},
				},
				// {
				// 	Label:   astikit.StrPtr("Fullscreen"),
				// 	Checked: astikit.BoolPtr(false),
				// 	Type:    astilectron.MenuItemTypeCheckbox,
				// 	OnClick: func(e astilectron.Event) (deleteListener bool) {
				// 		if *e.MenuItemOptions.Checked {
				// 			sui.awin.SetFullScreen(True)
				// 		} else {
				// 			sui.awin.SetFullScreen(False)
				// 		}
				// 		return
				// 	},
				// },
			},
		},
	})
	m.Create()
	sui.awin.Maximize()
}

func (sui *UI) Close() {
	sui.Asel.Close()
}

func (sui *UI) OpenDevTools() {
	sui.awin.OpenDevTools()
}

func (sui *UI) Wait() {
	sui.Asel.Wait()
}

func (sui *UI) createMSG(itime time.Duration, activescreen string, numstrike uint8) ([]byte, error) {
	type msg struct {
		Time   string `json:"time"`
		Screen string `json:"screen"`
		Strike uint8  `json:"strike"`
	}
	message, err := json.Marshal(msg{itime.String(), activescreen, numstrike})
	return message, err
}

func (sui *UI) UpdateUI(message []byte) {
	sui.awin.SendMessage(message)
}
