package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mktne "github.com/jax-b/ModulerKTNE/rpi_software/controller"
	"go.uber.org/zap"
)

func main() {
	cfg := &mktne.Config{}

	cfg.Shield.I2cBusNumber = 1
	cfg.Shield.Strike1PinNum = 22
	cfg.Shield.Strike2PinNum = 23
	cfg.Shield.ModInterruptPinNum = 17
	cfg.Shield.MfbStartPinNum = 27
	cfg.Shield.SevenSegAddress = 0x70

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	shieldctrl := mktne.NewShieldControl(sugar, cfg)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		shieldctrl.Close()
		os.Exit(0)
	}()

	shieldctrl.Run()

	sugar.Info("Hello World")
	sugar.Info("Waiting for button press")
	mfbchan := shieldctrl.RegisterMFBConsumer()
	mfbtime := <-mfbchan
	sugar.Infof("Got a button press with duration: %s", mfbtime)
	sugar.Info("Waiting for m2c signal")
	m2cchan := shieldctrl.RegisterM2CConsumer()
	<-m2cchan
	sugar.Info("Got m2c signal")

	sugar.Info("Attempting Close")
	shieldctrl.Close()
}
