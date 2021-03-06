package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
	"time"

	ctrl "github.com/jax-b/ModulerKTNE/rpi_software/controller"
)

func main() {
	boolPtr := flag.Bool("damon", false, "Run as damon (logs to file insted of stdout)")
	logger := ctrl.NewLogger(*boolPtr)
	ctrlr := ctrl.NewGameCtrlr(logger, false)
	reader := bufio.NewReader(os.Stdin)
	logger.Info("Installed Modules: ", ctrlr.GetInstalledModules())
	ctrlr.SetSerial("TQ74B01")
	ctrlr.SetNumBatteries(0)
	ctrlr.SetPorts(53)
	ctrlr.SetGameTime(4 * time.Minute)
	ctrlr.SetModSeed(0, 0xFFFF)
	ctrlr.SetStrikes(0, true)
	go ctrlr.Run()
	ctrlr.StartGame()
	ctrlr.UpdateModTime()
	logger.Info("press E to stop")
	for {
		input := make([]byte, 1)
		_, err := reader.Read(input)
		if err != nil {
			logger.Errorf("could not process input %v\n", input)
		}
		if strings.ToLower(string(input)) == "e" {
			logger.Info("Input Detected")
			break
		}
	}

	ctrlr.Close()
}
