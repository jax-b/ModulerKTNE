package main

import (
	"time"

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

	sugar.Info("Hello World")

	sugar.Info("Timer Beep")
	closech := make(chan bool)
	timetick := make(chan bool)

	shieldctrl.TimerBeep(closech, timetick)

	for i := 10; i > 0; i-- {
		sugar.Info(i)
		timetick <- true
		time.Sleep(time.Second)
	}

	closech <- true

	sugar.Info("Adding Strike")
	shieldctrl.AddStrike()
	time.Sleep(time.Millisecond * 500)
	sugar.Info("Playing Expload")
	shieldctrl.ExploadSound()
	time.Sleep(time.Millisecond * 500)
	sugar.Info("Playing Game Win")
	shieldctrl.GameWinSound()
	time.Sleep(time.Millisecond * 500)
	sugar.Info("Playing Module Solved")
	shieldctrl.ModSolvedSound()
	time.Sleep(time.Millisecond * 500)
	sugar.Info("Playing Module Solved")
	shieldctrl.NeedyWantSound()

	sugar.Info("Attempting Close")
	shieldctrl.Close()
}
