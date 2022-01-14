package controller

import (
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
					sgc.randomPopulate()
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
				os.exit(0)
				// syscall.Shutdown(0, 0)
			}
		case <-sgc.btnWatchStopCh:
			break
		}
	}
}
