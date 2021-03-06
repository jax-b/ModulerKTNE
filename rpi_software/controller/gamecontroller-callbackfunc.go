package controller

import (
	"fmt"
	"os"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
)

var (
	nextAnounce   time.Time     = time.Now().Add(5 * time.Second)
	nextModUpdate time.Time     = time.Now().Add(2 * time.Second)
	nextAudioTick time.Duration = 0
)

func (sgc *GameController) tmrCallbackFunction(tsb time.Time, tsa time.Time, stat mktnecf.Status) {
	fmt.Println(stat.Time)
	// Check for a boom
	if stat.Boom {
		sgc.log.Info("Timer Has Expired: BOOM!")
		sgc.GameOverBoom()
		err := sgc.updateNetworkIPC()
		if err != nil {
			sgc.log.Errorf("Failled to Announce the current status to external components (BOOM): %e", err)
		}
	}
	// Update external components about the current status of the game
	if nextAnounce.Before(time.Now()) {
		sgc.log.Info("Announcing Status")
		err := sgc.updateNetworkIPC()
		if err != nil {
			sgc.log.Errorf("Failled to Announce the current status to external components: %e", err)
		}
		nextAnounce = time.Now().Add(5 * time.Second)
	}
	// Resync the game clock to the modules
	if nextModUpdate.Before(time.Now()) {
		sgc.log.Info("Syncing Modules Time")
		err := sgc.UpdateModTime()
		if err != nil {
			sgc.log.Errorf("Failled to Sync the current gametime to the modules: %e", err)
		}
		nextModUpdate = time.Now().Add(2 * time.Second)
	}
	// Play Audio Tick
	if (nextAudioTick-stat.Time).Seconds() >= 1 || nextAudioTick == 0 {
		go func() { sgc.audio.tick <- true }()
		nextAudioTick = stat.Time
	}
	// Update Clock
	if stat.Time.Milliseconds()%10 == 0 {
		go sgc.rpishield.WriteTime(stat.Time)
	}
}

// MFB tracker
// Short Press will either start a new random game or stop the current game
// Long Press will shutdown the host os
func (sgc *GameController) buttonWatcher() {
	go func() {
		mfb := sgc.rpishield.RegisterMFBConsumer()
		log := sgc.log.Named("ButtonWatcher")
		for {
			select {
			case presstime := <-mfb:
				// wait for a button press
				if presstime > 50*time.Millisecond && presstime < 500*time.Millisecond { // Short Press
					log.Info("MFB Short Press Detected")
					if !sgc.game.comStat.Gamerun { // If the Game is not running, start a new game
						sgc.RandomPopulate()
						sgc.StartGame()
					} else { // If the Game is running, stop the game
						sgc.StopGame()
						err := sgc.updateNetworkIPC()
						if err != nil {
							log.Error(err)
						}
					}
				} else if presstime > 5*time.Second { // Long Press
					log.Info("MFB Long Press Detected")
					sgc.Close()
					os.Exit(0)
					// syscall.Shutdown(0, 0)
				}
			case <-sgc.btnWatchStopCh:
				return
			}
		}
	}()
}

// Handles the interrupt from the a modules updating its status in the game controller and updating strikes
func (sgc *GameController) m2cInterruptHandler() {
	go func() {
		interupt := sgc.rpishield.RegisterM2CConsumer()
		log := sgc.log.Named("M2CWatcher")
		for {
			select {
			case <-interupt:
				log.Info("Interrupt received")

				go func() {
					for index := range sgc.modules {
						if sgc.modules[index].present && !sgc.modules[index].solved {
							solvedStat, err := sgc.modules[index].mctrl.GetSolvedStatus()
							if err != nil {
								log.Error("Failed to get solved status", err)
							}

							if solvedStat < int8(int16(sgc.game.comStat.NumStrike)*-1) {
								log.Infof("Module %d Has A New Strike", index)
								sgc.AddStrike()
								sgc.UpdateModTime()
							} else if solvedStat > 0 {
								log.Infof("Module %d Has Been Solved", index)
								sgc.modules[index].solved = true
							}
						}
					}
				}()
			case <-sgc.interStopCh:
				return
			}
		}
	}()
}
