package main

import (
	"flag"

	ctrl "github.com/jax-b/ModulerKTNE/rpi_software/controller"
)

func main() {
	boolPtr := flag.Bool("damon", false, "Run as damon (logs to file insted of stdout)")
	logger := ctrl.NewLogger(*boolPtr)
	ctrlr := ctrl.NewGameCtrlr(logger)
	ctrlr.Run()
	defer ctrlr.Close()
}
