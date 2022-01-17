package main

import (
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
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
	var mods [10]*mktnecf.ModControl

	//Loop through the modules and create their control objects
	MCADDR := [10]byte{mktnecf.FRONT_MOD_1, mktnecf.FRONT_MOD_2, mktnecf.FRONT_MOD_3, mktnecf.FRONT_MOD_4, mktnecf.FRONT_MOD_5, mktnecf.BACK_MOD_1, mktnecf.BACK_MOD_2, mktnecf.BACK_MOD_3, mktnecf.BACK_MOD_4, mktnecf.BACK_MOD_5}
	for index, addr := range MCADDR {
		mods[index] = mktnecf.NewModControl(sugar, addr, int(cfg.Shield.I2cBusNumber))
	}

	for _, mctlr := range mods {
		modpresent := mctlr.TestIfPresent()
		if modpresent {
			sugar.Infof("Found Module: %x", mctlr.GetAddress())

			mctlr.SetGameSeed(0xFFFF)
			time.Sleep(time.Second * 5)
			mctlr.ClearGameSeed()
		}
	}
	logger.Info("Done")
}
