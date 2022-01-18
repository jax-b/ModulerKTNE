package controller

import (
	"os"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
)

var (
	nextAnounce   time.Time = time.Now().Add(5 * time.Second)
	nextModUpdate time.Time = time.Now().Add(5 * time.Second)
)

func (sgc *GameController) tmrCallbackFunction(tsb time.Time, tsa time.Time, stat mktnecf.Status) {
	// Check for a boom
	if stat.Boom {
		sgc.log.Info("Timer Has Expired: BOOM!")
		sgc.StopGame()
		sgc.game.comStat.Boom = true
		sgc.game.comStat.Win = false
		sgc.game.comStat.Gamerun = false
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
	if nextAnounce.Before(time.Now()) {
		sgc.log.Info("Syncing Modules Time")
		err := sgc.updateModTime()
		if err != nil {
			sgc.log.Errorf("Failled to Sync the current gametime to the modules: %e", err)
		}
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
			case presstimeint := <-mfb:
				// wait for a button press
				presstime := time.Duration(presstimeint) * time.Millisecond
				if presstime > 50*time.Millisecond && presstime < 200*time.Millisecond { // Short Press
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
				break
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
				for index := range sgc.modules {
					if sgc.modules[index].present && !sgc.modules[index].solved {
						solvedStat, err := sgc.modules[index].mctrl.GetSolvedStatus()
						if err != nil {
							log.Error("Failed to get solved status", err)
						}

						if solvedStat < int8(int16(sgc.game.comStat.NumStrike)*-1) {
							log.Infof("Module %d Has A New Strike", index)
							sgc.AddStrike()
						} else if solvedStat > 0 {
							log.Infof("Module %d Has Been Solved", index)
							sgc.modules[index].solved = true
						}
					}
				}
			case <-sgc.interStopCh:
				return
			}
		}
	}()
}
