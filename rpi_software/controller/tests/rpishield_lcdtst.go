package main

import (
	mktne "github.com/jax-b/ModulerKTNE/rpi_software/controller"
	"go.uber.org/zap"
	"time"
)

func main() {
	cfg := &mktne.Config{}

	cfg.Shield.I2cBusNumber=1
	cfg.Shield.Strike1PinNum= 22
	cfg.Shield.Strike2PinNum=  23
	cfg.Shield.ModInterruptPinNum= 17
	cfg.Shield.MfbStartPinNum     =27
	cfg.Shield.SevenSegAddress= 0x70

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	shieldctrl := mktne.NewShieldControl(sugar, cfg)

	sugar.Info("Hello World")
	sugar.Info("Writing Idle")
	shieldctrl.WriteIdle()
	time.Sleep(time.Second * 1)
	sugar.Info("Writing Time")
	shieldctrl.WriteTime(time.Minute + time.Second * 30)
	time.Sleep(time.Second * 1)

	sugar.Info("Attempting Close")
	shieldctrl.Close()

}