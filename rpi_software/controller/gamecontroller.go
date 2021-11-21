package controller

import (
	"github.com/jax-b/ModulerKTNE/rpi_software/controller"
)

type GameController struct {
	sidePanels [4]*controller.SideControl
	modules    [10]struct {
		mctrl   *controller.ModControl
		present bool
	}
	rpishield  *controller.ShieldControl
	ipc        *controller.InterProcessCom
	mnetc      *controller.MultiCastCountdown
	cfg        *controller.Config
	gameTime   uint16
	gameStopCh chan bool
	numStrike  uint8
}

// Gets the current game time
func (self *GameController) GetTime() uint16 {
	return &self.gameTime
}

// Sets the game time to the given time
func (self *GameController) SetTime(time uint16) {
	for mod := range modules {
		if mod.present {
			mod.mctrl.SetTime(time)
		}
	}
	self.gameTime = time
}
func (self *GameController) StartGame() error {
	// for all the modules that are present, start the game
	self.scanModules()
	for mod := range modules {
		if mod.present {
			mod.mctrl.StartGame()
		}
	}
	self.gameStopCh = make(chan bool)
	go self.timer(self.gameStopCh)
	return nil
}
func (self *GameController) StopGame() error {
	// for all the modules that are present, stop the game
	for mod := range modules {
		if mod.present {
			mod.mctrl.StopGame()
		}
	}

	self.gameStopCh <- true
	return nil
}
func (self *GameController) timer(StopCh chan bool) {
	ticker := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-StopCh:
			return
		case <-ticker:
			self.gameTime--
		}
		if self.gameTime == 0 {
			self.timerRunOut()
			return
		}
	}
}
func (self *GameController) scanModules() {
	for mod := range modules {
		mod.present = mod.TestIfPresent()
	}
}
func (self *GameController) timerRunOut() {
	self.StopGame()
	self.mnetc.SendTimeStrike(0, self.numStrike, true, false)
	self.ipc.TimerRunOut()
}
