package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
	"github.com/jax-b/ModulerKTNE/rpi_software/testnetdata"
	"go.uber.org/zap"
)

var (
	logger      *zap.SugaredLogger
	mcastCount  *mktnecf.MultiCastCountdown
	stat        *mktnecf.Status
	nextAnounce time.Time
)

func startupInstructions() {
	fmt.Println("Press B to toggle boom")
	fmt.Println("Press T to toggle timmer run")
	fmt.Println("Press R to reset the timer and strikes")
	fmt.Println("Press S to add a strike")
	fmt.Println("Press W to toggle win")
	fmt.Println("Press E to exit")
}
func consolCMD() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input := make([]byte, 1)
		_, err := reader.Read(input)
		if err != nil {
			fmt.Printf("could not process input %v\n", input)
		}
		switch strings.ToLower(string(input)) {
		case "b":
			stat.Boom = !stat.Boom
			logger.Info("Toggled boom: ", stat.Boom)
			mcastCount.SendStatus(stat)
		case "w":
			stat.Win = !stat.Win
			if stat.Win {
				stat.Gamerun = false
				stat.Boom = false
			}
			logger.Info("Toggled win: ", stat.Win)
			mcastCount.SendStatus(stat)
		case "s":
			stat.NumStrike++
			logger.Info("Added strike: ", stat.NumStrike)
			mcastCount.SendStatus(stat)
		case "t":
			stat.Gamerun = !stat.Gamerun
			nextAnounce = time.Now().Add(5 * time.Second)
			logger.Info("Toggled timer: ", stat.Gamerun)
			mcastCount.SendStatus(stat)
		case "r":
			stat.NumStrike = 0
			stat.Time = 90 * time.Second
			stat.Boom = false
			stat.Win = false
			logger.Info("reset time to 90 seconds and strikes to 0")
			mcastCount.SendStatus(stat)
		case "e":
			stat.NumStrike = 0
			stat.Time = 1 * time.Hour
			stat.Boom = false
			stat.Win = false
			stat.Gamerun = false
			mcastCount.SendStatus(stat)
			os.Exit(0)
		}
	}
}

func callbackfunc(tsb time.Time, tsa time.Time, stat mktnecf.Status) {
	logger.Debug("Timeleft: ", stat.Time)
	if stat.Boom {
		logger.Info("Boom!")
		mcastCount.SendStatus(&stat)
	}
	if nextAnounce.Before(time.Now()) {
		logger.Infof("Announce: Timeleft: %s, Strikes: %d, Boom: %t, Win: %t, trunning: %t, strikerate: %f ", stat.Time.String(), stat.NumStrike, stat.Boom, stat.Win, stat.Gamerun, stat.Strikereductionrate)
		mcastCount.SendStatus(&stat)
		nextAnounce = time.Now().Add(5 * time.Second)
	}
}

func main() {
	startupInstructions()
	logger = testnetdata.NewLogger()
	config := testnetdata.NewConfig(logger)
	config.Load()
	logger.Info("Config:", config)
	stat = &mktnecf.Status{
		Time:                15 * time.Second,
		NumStrike:           0,
		Boom:                false,
		Win:                 false,
		Gamerun:             true,
		Strikereductionrate: float32(0.25),
	}

	var err error
	mcastCount, err = mktnecf.NewMultiCastCountdown(logger, config.Network.MultiCastIP, config.Network.MultiCastPort)
	if err != nil {
		logger.Error("Error:", err)
		os.Exit(1)
	}
	defer mcastCount.Close()
	tmr := mktnecf.NewGameTimer(logger, stat)
	tmr.AddCallbackFunction(callbackfunc)
	defer tmr.Close()
	go consolCMD()

	logger.Infof("Announce: Timeleft: %s, Strikes: %d, Boom: %t, Win: %t, trunning: %t, strikerate: %f ", stat.Time.String(), stat.NumStrike, stat.Boom, stat.Win, stat.Gamerun, stat.Strikereductionrate)
	nextAnounce = time.Now().Add(5 * time.Second)
	mcastCount.SendStatus(stat)

	tmr.Run()
}
