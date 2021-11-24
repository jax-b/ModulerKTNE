package controller

import (
	"time"

	"go.uber.org/zap"
)

type module struct {
	mctrl   *ModControl
	present bool
}
type multicast struct {
	useMulti bool
	mnetc    *MultiCastCountdown
}
type indicator struct {
	lit   bool
	label [3]rune
}
type gameinfo struct {
	time       uint32
	numStrike  int8 //Works just like mod 1 is solved anything negative is a strike
	run        bool
	strikerate float32
	indicator  []indicator
	port       []byte
}
type GameController struct {
	sidePanels [4]*SideControl
	modules    [10]module
	multicast  multicast
	game       gameinfo
	rpishield  *ShieldControl
	ipc        *InterProcessCom
	cfg        *Config
	gameStopCh chan bool
	log        *zap.SugaredLogger
}

func NewGameCtrlr() *GameController {
	// Set up the program logger
	zapl, _ := zap.NewProduction()
	logger := zapl.Sugar()
	logger = logger.Named("MKTNE")

	// Load the configuration
	cfg := NewConfig(logger)
	cfg.Load()

	// Try to open the RPI Shield Control
	rpis := NewShieldControl(logger, cfg)

	// Create the GameController Object
	gc := &GameController{
		rpishield: rpis,
		cfg:       cfg,
		log:       logger,
		game: gameinfo{
			run:        false,
			numStrike:  0,
			time:       0,
			strikerate: 0.25,
		},
	}

	// Create the inter process communicator object
	gc.ipc = NewIPC(logger, gc)

	//Loop through the side panels and create their control objects
	SPADDR := [4]byte{TOP_PANEL, RIGHT_PANEL, BOTTOM_PANEL, LEFT_PANEL}
	for i := range gc.sidePanels {
		gc.sidePanels[i] = NewSideControl(gc.log, SPADDR[i], int(gc.cfg.Shield.I2cBusNumber))
	}

	//Loop through the modules and create their control objects
	MCADDR := [10]byte{FRONT_MOD_1, FRONT_MOD_2, FRONT_MOD_3, FRONT_MOD_4, FRONT_MOD_5, BACK_MOD_1, BACK_MOD_2, BACK_MOD_3, BACK_MOD_4, BACK_MOD_5}
	for i := range gc.modules {
		gc.modules[i].mctrl = NewModControl(gc.log, MCADDR[i], int(gc.cfg.Shield.I2cBusNumber))
	}

	// Check if multicast is enabled then create its object
	if gc.cfg.Network.UseMulticast {
		gc.multicast.useMulti = true
		var err error
		gc.multicast.mnetc, err = NewMultiCastCountdown(gc.log, gc.cfg)
		if err != nil {
			gc.log.Error("Failed to create multicast countdown, proceeding without multicast", err)
			gc.multicast.useMulti = false
		}
	} else {
		gc.multicast.useMulti = false
	}

	return gc
}

// Safe Shutdown of all components
func (self *GameController) Close() {
	// flush the logger
	self.log.Sync()
	// Close all of the modules
	for i := range self.modules {
		self.modules[i].mctrl.Close()
	}
	// Close the RPI Shield
	for i := range self.sidePanels {
		self.sidePanels[i].Close()
	}
	// Close the shield
	self.rpishield.Close()

	// Close multicast if used
	if self.multicast.useMulti {
		self.multicast.mnetc.Close()
	}

	//Close the IPC
	self.ipc.Close()
}

// Gets the current game time
func (self *GameController) GetGameTime() uint32 {
	return self.game.time
}

// Sets the game time to the given time
func (self *GameController) SetGameTime(time uint32) error {
	self.game.time = time
	self.UpdateModTime()
	return nil
}

// Updates the time on a module to the current game time
func (self *GameController) UpdateModTime() error {
	for i := range self.modules {
		if self.modules[i].present {
			err := self.modules[i].mctrl.SyncGameTime(self.game.time)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Gets the current amount of strikes
func (self *GameController) GetStrikes() int8 {
	return self.game.numStrike
}

// Set the number of strikes
func (self *GameController) SetStrikes(strikes int8) error {
	for i := range self.modules {
		if self.modules[i].present {
			err := self.modules[i].mctrl.SetSolvedStatus(strikes * -1)
			if err != nil {
				return err
			}
		}
	}
	self.game.numStrike = strikes
	return nil
}

// Get the srike reduction rate
func (self *GameController) GetStrikeRate() float32 {
	return self.game.strikerate
}

// Set the strike reduction rate
func (self *GameController) SetStrikeRate(rate float32) error {
	for i := range self.modules {
		if self.modules[i].present {
			err := self.modules[i].mctrl.SetStrikeReductionRate(rate)
			if err != nil {
				return err
			}
		}
	}
	self.game.strikerate = rate
	return nil
}

// Starts the game
func (self *GameController) StartGame() error {
	// for all the modules that are present, start the game
	self.scanModules()
	for i := range self.modules {
		if self.modules[i].present {
			self.modules[i].mctrl.StartGame()
		}
	}
	self.gameStopCh = make(chan bool)
	self.game.run = true
	go self.timer(self.gameStopCh)
	return nil
}

// Stops the game
func (self *GameController) StopGame() error {
	if self.game.run {
		// for all the modules that are present, stop the game
		for i := range self.modules {
			if self.modules[i].present {
				self.modules[i].mctrl.StopGame()
			}
		}

		self.gameStopCh <- true
	}
	return nil
}

// ******** Will have to have a bunch of work done on it to include all of the necessary functions
// this function is the timekeeper for the game
func (self *GameController) timer(StopCh chan bool) {
	ticker := time.NewTicker(time.Millisecond)
	countTicker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-StopCh:
			return
		case <-ticker.C:
			// Need to add reduction rate
			self.game.time--
		case <-countTicker.C:
			go self.UpdateModTime()
			go self.ipc.SyncStatus(self.game.time, self.game.numStrike, false, false)
			if self.multicast.useMulti {
				go self.multicast.mnetc.SendStatus(self.game.time, self.game.numStrike, false, false)
			}
		}
		if self.game.time == 0 {
			self.timerRunOut()
			return
		}
	}
}

// If the timer is to run out here is how we handle it
func (self *GameController) timerRunOut() {
	self.StopGame()
	if self.multicast.useMulti {
		self.multicast.mnetc.SendStatus(0, self.game.numStrike, true, false)
	}
	self.ipc.SyncStatus(0, self.game.numStrike, true, false)
}

// Polls all possible module addresses and sees if something is their. Updates the class variables
func (self *GameController) scanModules() {
	for i := range self.modules {
		self.modules[i].present = self.modules[i].mctrl.TestIfPresent()
	}
}
