package main

import (
	"github.com/jax-b/ModulerKTNE/rpi_software/netdisplay"
	"go.uber.org/zap"
)

var (
	zaplogger *zap.SugaredLogger
)

func main() {
	// Create a new logger
	zaplogger = netdisplay.NewLogger()
	netdisp := netdisplay.NewNetDisplay(zaplogger)
	netdisp.Run()
	netdisp.UI.OpenDevTools()
	// Blocking pattern
	netdisp.UI.Asel.Wait()
	netdisp.Close()
}
