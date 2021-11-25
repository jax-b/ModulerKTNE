package controller

import (
	"math/rand"
	"syscall"
	"time"

	"github.com/jax-b/ModulerKTNE/rpi_software/util"
	"go.uber.org/zap"
)

var (
	VALID_INDICATORS       []string = []string{"SND", "CLR", "CAR", "IND", "FRQ", "SIG", "NSA", "MSA", "TRN", "BOB", "FRK"}
	VALID_PORT_ID          []string = []string{"DVI", "PAR", "PS2", "RJ4", "SER", "RCA"}
	VALID_PORTS_TRANSLATED []byte   = []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
)

const (
	GAMEPLAYMAXNUMPORT    = 6
	GAMEPLAYMAXTINDICATOR = 12
)

type module struct {
	mctrl   *ModControl
	present bool
}
type multicast struct {
	useMulti bool
	mnetc    *MultiCastCountdown
}
type Indicator struct {
	lit   bool    `json:lit`
	label [3]rune `json:label`
}
type gameinfo struct {
	time       uint32
	numStrike  int8 //Works just like mod 1 is solved anything negative is a strike
	run        bool
	strikerate float32
	indicators []Indicator
	port       []uint8
	serialnum  [8]rune
	numbat     int
}
type GameController struct {
	sidePanels     [4]*SideControl
	modules        [10]module
	multicast      multicast
	game           gameinfo
	rpishield      *ShieldControl
	ipc            *InterProcessCom
	cfg            *Config
	timerStopCh    chan bool
	btnWatchStopCh chan bool
	log            *zap.SugaredLogger
	rnd            *rand.Rand
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
		rpishield:      rpis,
		cfg:            cfg,
		log:            logger,
		btnWatchStopCh: make(chan bool),
		timerStopCh:    make(chan bool),
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

	// Initalize RNG
	var src util.CryptoSource
	rnd := rand.New(src)
	gc.rnd = rnd

	go gc.buttonWatcher()

	return gc
}

// Safe Shutdown of all components
func (self *GameController) Close() {
	go func() { self.timerStopCh <- true }()
	self.btnWatchStopCh <- true
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

// Adds a indicator to the list
func (self *GameController) AddIndicator(indi Indicator) {
	if len(self.game.indicators) > GAMEPLAYMAXTINDICATOR {
		self.game.indicators[len(self.game.indicators)] = indi
	} else {
		self.game.indicators = append(self.game.indicators, indi)
	}
	if indi.lit {
		for i := range self.modules {
			if self.modules[i].present {
				self.modules[i].mctrl.SetGameLitIndicator(indi.label)
			}
		}
	}
}

// Gets the currently configured indicators
func (self *GameController) GetIndicators() []Indicator {
	return self.game.indicators
}

// Clears out the current indicators
func (self *GameController) ClearIndicators() {
	for i := range self.modules {
		if self.modules[i].present {
			self.modules[i].mctrl.ClearGameLitIndicator()
		}
	}
	self.game.indicators = make([]Indicator, 0)
}

// Adds a port to the list
func (self *GameController) AddPort(port byte) error {
	self.game.port = append(self.game.port, port)
	for i := range self.modules {
		if self.modules[i].present {
			err := self.modules[i].mctrl.SetGamePortID(port)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// returns all of the ports that are configured for the game
func (self *GameController) GetPorts() []byte {
	return self.game.port
}

// clears all of the ports that are configured for the game
func (self *GameController) ClearPorts() error {
	self.game.port = make([]byte, 0)
	for i := range self.modules {
		if self.modules[i].present {
			err := self.modules[i].mctrl.ClearGamePortIDS()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Sets the current game serial number
func (self *GameController) SetSerial(serial string) error {
	for i := range serial {
		if i > len(self.game.serialnum) {
			break
		}
		self.game.serialnum[i] = rune(serial[i])
	}
	for i := range self.modules {
		if self.modules[i].present {
			err := self.modules[i].mctrl.SetGameSerialNumber(self.game.serialnum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *GameController) GetSerial() string {
	return string(self.game.serialnum[0:])
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
	self.timerStopCh = make(chan bool)
	self.game.run = true
	go self.timer(self.timerStopCh)
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

		self.timerStopCh <- true
	}
	return nil
}

// Update the specified module with all of the game values
func (self *GameController) ModFullUpdate(modnum int) {
	var litindi [][3]rune
	for i := range self.game.indicators {
		if self.game.indicators[i].lit {
			litindi = append(litindi, self.game.indicators[i].label)
		}
	}
	self.modules[modnum].mctrl.SetupAllGameData(
		self.game.serialnum,
		litindi,
		uint8(self.game.numbat),
		self.game.port,
	)
	self.modules[modnum].mctrl.SetStrikeReductionRate(self.game.strikerate)
	self.modules[modnum].mctrl.SetSolvedStatus(self.game.numStrike)
}

// MFB tracker
func (self *GameController) buttonWatcher() {
	mfb := self.rpishield.RegisterMFBConsumer()
	for {
		select {
		case presstimeint := <-mfb:
			// wait for a button press
			presstime := time.Duration(presstimeint) * time.Millisecond
			if presstime > 50 {
				if !self.game.run {
					self.randomPopulate()
					self.StartGame()
				} else {
					self.timerRunOut()
				}
			} else if presstime > 5*time.Second {
				self.Close()
				syscall.Shutdown(0, 0)
			}
		case <-self.btnWatchStopCh:
			break
		}
	}
}

func (self *GameController) randomPopulate() {
	// Serial number generation
	serialLen := self.rnd.Intn(8-6) + 6
	serial := make([]byte, serialLen)
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range serial {
		serial[i] = charset[self.rnd.Intn(len(charset))]
	}
	self.SetSerial(string(serial))

	// Indicator Generation
	numIndicator := self.rnd.Intn(GAMEPLAYMAXTINDICATOR)
	self.ClearIndicators()
	for i := 0; i < numIndicator; i++ {
		indilbl := VALID_INDICATORS[self.rnd.Intn(len(VALID_INDICATORS))]
		indilit := self.rnd.Intn(2) == 1
		var indilblrn [3]rune
		for i := range indilblrn {
			indilblrn[i] = rune(indilbl[i])
		}
		indi := Indicator{
			label: indilblrn,
			lit:   indilit,
		}
		self.AddIndicator(indi)
	}

	// Port Generation
	numPort := self.rnd.Intn(GAMEPLAYMAXNUMPORT)
	self.ClearPorts()
	for i := 0; i < numPort; i++ {
		port := VALID_PORTS_TRANSLATED[self.rnd.Intn(len(VALID_PORT_ID))]
		self.AddPort(port)
	}
}

// this function is the timekeeper for the game
func (self *GameController) timer(StopCh chan bool) {
	ticker := time.NewTicker(time.Millisecond)
	countTicker := time.NewTicker(time.Second * 30)
	extratick := 0
	for {
		select {
		case <-StopCh:
			return
		case <-ticker.C:
			// Need to add reduction rate
			self.game.time--
			if self.game.numStrike < 0 {
				everyrate := int((1 / self.game.strikerate) / (-1 * float32(self.game.numStrike)))
				if extratick >= everyrate {
					self.game.time--
					extratick = 0
				} else {
					extratick++
				}
			}
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
		laststate := self.modules[i].present
		self.modules[i].present = self.modules[i].mctrl.TestIfPresent()
		if laststate != self.modules[i].present && self.modules[i].present {
			self.ModFullUpdate(i)
		}
	}
}
