package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jax-b/ModulerKTNE/rpi_software/testnetdata"
	"go.uber.org/zap"
)

var (
	boom          = false
	win           = false
	strikes       = 0
	logger        *zap.SugaredLogger
	mcastCount    *testnetdata.MultiCastCountdown
	countdowntime time.Duration
)

func startupInstructions() {
	fmt.Println("Press B to toggle boom")
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
			boom = !boom
			logger.Info("Toggled boom: ", boom)
			mcastCount.SendStatus(uint32(countdowntime.Milliseconds()), int8(strikes), boom, win)
		case "w":
			win = !win
			logger.Info("Toggled win: ", win)
			mcastCount.SendStatus(uint32(countdowntime.Milliseconds()), int8(strikes), boom, win)
		case "s":
			strikes++
			logger.Info("Added strike: ", strikes)
			mcastCount.SendStatus(uint32(countdowntime.Milliseconds()), int8(strikes), boom, win)
		case "r":
			strikes = 0
			countdowntime = 10 * time.Minute
			logger.Info("reset time an strikes: ", strikes)
			mcastCount.SendStatus(uint32(countdowntime.Milliseconds()), int8(strikes), boom, win)
		case "e":
			os.Exit(0)
		}
	}
}
func gTimer() {
	timeticker := time.NewTicker(time.Millisecond * 1)
	timeanounceticker := time.NewTicker(time.Second * 5)
	extratick := 0
	for {
		select {
		case <-timeticker.C:
			countdowntime = countdowntime - time.Millisecond
			if strikes < 0 {
				everyrate := int((1 / 0.25) / float32(strikes))
				if extratick >= everyrate {
					countdowntime = countdowntime - time.Millisecond
					extratick = 0
				} else {
					extratick++
				}
			}
		case <-timeanounceticker.C:
			inttime := uint32(countdowntime.Milliseconds())
			mcastCount.SendStatus(inttime, int8(strikes), boom, win)
			logger.Infof("Announce: Timeleft: %d, Strikes: %d, Boom: %t, Win: %t", inttime, strikes, boom, win)
		}
	}
}

func main() {
	startupInstructions()
	logger = testnetdata.NewLogger()
	config := testnetdata.NewConfig(logger)
	config.Load()
	logger.Info("Config:", config)
	var err error
	mcastCount, err = testnetdata.NewMultiCastCountdown(logger, config)
	if err != nil {
		logger.Error("Error:", err)
		os.Exit(1)
	}
	defer mcastCount.Close()

	countdowntime = 10 * time.Minute
	go consolCMD()
	gTimer()
}
