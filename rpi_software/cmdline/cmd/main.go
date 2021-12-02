package main

import (
	"flag"
	"fmt"

	cline "github.com/jax-b/ModulerKTNE/rpi_software/cmdline"
)

func main() {
	nop := true
	customcommandPTR := flag.String("cc", "", "Sends a custom command to the server")
	flag.Parse()

	ipc := cline.NewIPC()
	defer ipc.Close()
	if *customcommandPTR != "" {
		_, err := ipc.SendCustom(*customcommandPTR)
		if err != nil {
			fmt.Println(err)
		}
		nop = false
	}

	if !nop {
		flag.PrintDefaults()
	}

}
