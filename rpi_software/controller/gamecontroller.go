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
	solved  bool
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
	numstrike  int8 //Works just like mod 1 is solved anything negative is a strike
	run        bool
	strikerate float32
	indicators []Indicator
	port       []uint8
	serialnum  [8]rune
	numbat     int
	maxstrike  int8
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
	interStopCh    chan bool
	solvedStopCh   chan bool
	log            *zap.SugaredLogger
	rnd            *rand.Rand
}

func NewGameCtrlr(runAsDamon bool) *GameController {
	// Set up the program logger
	logger := NewLogger(runAsDamon)

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
		interStopCh:    make(chan bool),
		solvedStopCh:   make(chan bool),
		game: gameinfo{
			run:        false,
			numstrike:  0,
			maxstrike:  2,
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

	return gc
}
func (sgc *GameController) Run() {
	sgc.ipc.Run()
	sgc.rpishield.Run()
	go sgc.buttonWatcher()
	go sgc.interruptHandler()
	go sgc.solvedCheck()
}

// Safe Shutdown of all components
func (sgc *GameController) Close() {
	go func() { sgc.timerStopCh <- true }()
	sgc.btnWatchStopCh <- true
	sgc.interStopCh <- true
	sgc.solvedStopCh <- true
	// flush the logger
	sgc.log.Sync()
	// Close all of the modules
	for i := range sgc.modules {
		sgc.modules[i].mctrl.Close()
	}
	// Close the RPI Shield
	for i := range sgc.sidePanels {
		sgc.sidePanels[i].Close()
	}
	// Close the shield
	sgc.rpishield.Close()

	// Close multicast if used
	if sgc.multicast.useMulti {
		sgc.multicast.mnetc.SendStatus(0, 0, false, false, false, sgc.game.strikerate)
		sgc.multicast.mnetc.Close()
	}
	//Close the IPC
	sgc.ipc.Close()
}

// Checks to see if each module is solved and if all of them are solved then trigger the win condition
func (sgc *GameController) solvedCheck() {
	for {
		select {
		case <-sgc.solvedStopCh:
			return
		default:
			notSolved := false
			for i := range sgc.modules {
				if sgc.modules[i].present && !sgc.modules[i].solved {
					notSolved = true
				}
			}
			if notSolved {
				sgc.StopGame()
				sgc.ipc.SyncStatus(sgc.game.time, sgc.game.numstrike, false, true)
				if sgc.multicast.useMulti {
					sgc.multicast.mnetc.SendStatus(sgc.game.time, sgc.game.numstrike, false, true, false, sgc.game.strikerate)
				}
			}
		}
	}
}

// Handles the interrupt from the a modules updating its status in the game controller and updating strikes
func (sgc *GameController) interruptHandler() {
	interupt := sgc.rpishield.RegisterM2CConsumer()
	for {
		select {
		case <-interupt:
			sgc.log.Info("Interrupt received")
			for i := range sgc.modules {
				if sgc.modules[i].present && !sgc.modules[i].solved {
					solvedStat, err := sgc.modules[i].mctrl.GetSolvedStatus()
					if err != nil {
						sgc.log.Error("Failed to get solved status", err)
					}
					if solvedStat < sgc.game.numstrike {
						sgc.AddStrike()
					} else if solvedStat > 0 {
						sgc.modules[i].solved = true
					}
				}
			}
		case <-sgc.interStopCh:
			return
		}
	}
}

// Gets the current game time
func (sgc *GameController) GetGameTime() uint32 {
	return sgc.game.time
}

// Sets the game time to the given time
func (sgc *GameController) SetGameTime(time uint32) error {
	sgc.game.time = time
	sgc.UpdateModTime()
	return nil
}

// Updates the time on a module to the current game time
func (sgc *GameController) UpdateModTime() error {
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SyncGameTime(sgc.game.time)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Gets the current amount of strikes
func (sgc *GameController) GetStrikes() int8 {
	return sgc.game.numstrike
}

// Set the number of strikes
func (sgc *GameController) SetStrikes(strikes int8) error {
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetSolvedStatus(strikes)
			if err != nil {
				return err
			}
		}
	}
	sgc.game.numstrike = strikes
	return nil
}

func (sgc *GameController) AddStrike() {
	sgc.game.numstrike--
	if sgc.game.numstrike > sgc.game.maxstrike {
		if sgc.multicast.useMulti {
			sgc.multicast.mnetc.SendStatus(sgc.game.time, sgc.game.numstrike, true, false, sgc.game.run, sgc.game.strikerate)
		}
		sgc.ipc.SyncStatus(sgc.game.time, sgc.game.numstrike, true, false)
	} else {
		for i := range sgc.modules {
			if sgc.modules[i].present && !sgc.modules[i].solved {
				err := sgc.modules[i].mctrl.SetSolvedStatus(sgc.game.numstrike)
				if err != nil {
					sgc.log.Error("Failed to set solved status", err)
				}
			}
		}
		sgc.rpishield.AddStrike()
		if sgc.multicast.useMulti {
			sgc.multicast.mnetc.SendStatus(sgc.game.time, sgc.game.numstrike, false, false, sgc.game.run, sgc.game.strikerate)
		}
		sgc.ipc.SyncStatus(sgc.game.time, sgc.game.numstrike, false, false)
	}
}

// Get the srike reduction rate
func (sgc *GameController) GetStrikeRate() float32 {
	return sgc.game.strikerate
}

// Set the strike reduction rate
func (sgc *GameController) SetStrikeRate(rate float32) error {
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetStrikeReductionRate(rate)
			if err != nil {
				return err
			}
		}
	}
	sgc.game.strikerate = rate
	return nil
}

// Adds a indicator to the list
func (sgc *GameController) AddIndicator(indi Indicator) {
	if len(sgc.game.indicators) > GAMEPLAYMAXTINDICATOR {
		sgc.game.indicators[len(sgc.game.indicators)] = indi
	} else {
		sgc.game.indicators = append(sgc.game.indicators, indi)
	}
	if indi.lit {
		for i := range sgc.modules {
			if sgc.modules[i].present {
				sgc.modules[i].mctrl.SetGameLitIndicator(indi.label)
			}
		}
	}
}

// Gets the currently configured indicators
func (sgc *GameController) GetIndicators() []Indicator {
	return sgc.game.indicators
}

// Clears out the current indicators
func (sgc *GameController) ClearIndicators() {
	for i := range sgc.modules {
		if sgc.modules[i].present {
			sgc.modules[i].mctrl.ClearGameLitIndicator()
		}
	}
	sgc.game.indicators = make([]Indicator, 0)
}

// Adds a port to the list
func (sgc *GameController) AddPort(port byte) error {
	sgc.game.port = append(sgc.game.port, port)
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetGamePortID(port)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// returns all of the ports that are configured for the game
func (sgc *GameController) GetPorts() []byte {
	return sgc.game.port
}

// clears all of the ports that are configured for the game
func (sgc *GameController) ClearPorts() error {
	sgc.game.port = make([]byte, 0)
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.ClearGamePortIDS()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Sets the current game serial number
func (sgc *GameController) SetSerial(serial string) error {
	for i := range serial {
		if i > len(sgc.game.serialnum) {
			break
		}
		sgc.game.serialnum[i] = rune(serial[i])
	}
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetGameSerialNumber(sgc.game.serialnum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sgc *GameController) GetSerial() string {
	return string(sgc.game.serialnum[0:])
}

// Starts the game
func (sgc *GameController) StartGame() error {
	// for all the modules that are present, start the game
	sgc.scanModules()
	for i := range sgc.modules {
		if sgc.modules[i].present {
			sgc.modules[i].solved = false
			sgc.modules[i].mctrl.StartGame()
		}
	}
	sgc.timerStopCh = make(chan bool)
	sgc.game.run = true
	go sgc.timer(sgc.timerStopCh)
	return nil
}

// Stops the game
func (sgc *GameController) StopGame() error {
	if sgc.game.run {
		// for all the modules that are present, stop the game
		for i := range sgc.modules {
			if sgc.modules[i].present {
				sgc.modules[i].mctrl.StopGame()
			}
		}

		sgc.timerStopCh <- true
	}
	return nil
}

// Update the specified module with all of the game values
func (sgc *GameController) ModFullUpdate(modnum int) {
	var litindi [][3]rune
	for i := range sgc.game.indicators {
		if sgc.game.indicators[i].lit {
			litindi = append(litindi, sgc.game.indicators[i].label)
		}
	}
	sgc.modules[modnum].mctrl.SetupAllGameData(
		sgc.game.serialnum,
		litindi,
		uint8(sgc.game.numbat),
		sgc.game.port,
	)
	sgc.modules[modnum].mctrl.SetStrikeReductionRate(sgc.game.strikerate)
	sgc.modules[modnum].mctrl.SetSolvedStatus(sgc.game.numstrike)
}

// MFB tracker
func (sgc *GameController) buttonWatcher() {
	mfb := sgc.rpishield.RegisterMFBConsumer()
	for {
		select {
		case presstimeint := <-mfb:
			// wait for a button press
			presstime := time.Duration(presstimeint) * time.Millisecond
			if presstime > 50 {
				if !sgc.game.run {
					sgc.randomPopulate()
					sgc.StartGame()
				} else {
					sgc.timerRunOut()
				}
			} else if presstime > 5*time.Second {
				sgc.Close()
				syscall.Shutdown(0, 0)
			}
		case <-sgc.btnWatchStopCh:
			break
		}
	}
}

// populates each module with a random game
func (sgc *GameController) randomPopulate() {
	// Serial number generation
	serialLen := sgc.rnd.Intn(8-6) + 6
	serial := make([]byte, serialLen)
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range serial {
		serial[i] = charset[sgc.rnd.Intn(len(charset))]
	}
	sgc.SetSerial(string(serial))

	// Indicator Generation
	numIndicator := sgc.rnd.Intn(GAMEPLAYMAXTINDICATOR)
	sgc.ClearIndicators()
	for i := 0; i < numIndicator; i++ {
		indilbl := VALID_INDICATORS[sgc.rnd.Intn(len(VALID_INDICATORS))]
		indilit := sgc.rnd.Intn(2) == 1
		var indilblrn [3]rune
		for i := range indilblrn {
			indilblrn[i] = rune(indilbl[i])
		}
		indi := Indicator{
			label: indilblrn,
			lit:   indilit,
		}
		sgc.AddIndicator(indi)
	}

	// Port Generation
	numPort := sgc.rnd.Intn(GAMEPLAYMAXNUMPORT)
	sgc.ClearPorts()
	for i := 0; i < numPort; i++ {
		port := VALID_PORTS_TRANSLATED[sgc.rnd.Intn(len(VALID_PORT_ID))]
		sgc.AddPort(port)
	}
}

// this function is the timekeeper for the game
// it will stop the game when the time runs out
func (sgc *GameController) timer(StopCh chan bool) {
	ticker := time.NewTicker(time.Millisecond)
	countTicker := time.NewTicker(time.Second * 10)
	extratick := 0
	for {
		select {
		case <-StopCh:
			return
		case <-ticker.C:
			// Need to add reduction rate
			sgc.game.time--
			if sgc.game.numstrike < 0 {
				everyrate := int((1 / sgc.game.strikerate) / (-1 * float32(sgc.game.numstrike)))
				if extratick >= everyrate {
					sgc.game.time--
					extratick = 0
				} else {
					extratick++
				}
			}
		case <-countTicker.C:
			go sgc.UpdateModTime()
			go sgc.ipc.SyncStatus(sgc.game.time, sgc.game.numstrike, false, false)
			if sgc.multicast.useMulti {
				go sgc.multicast.mnetc.SendStatus(sgc.game.time, sgc.game.numstrike, false, false, sgc.game.run, sgc.game.strikerate)
			}
		}
		if sgc.game.time == 0 {
			sgc.timerRunOut()
			return
		}
	}
}

// If the timer is to run out here is how we handle it
func (sgc *GameController) timerRunOut() {
	sgc.StopGame()
	if sgc.multicast.useMulti {
		sgc.multicast.mnetc.SendStatus(0, sgc.game.numstrike, true, false, false, sgc.game.strikerate)
	}
	sgc.ipc.SyncStatus(0, sgc.game.numstrike, true, false)
}

// Polls all possible module addresses and sees if something is their. Updates the class variables
func (sgc *GameController) scanModules() {
	for i := range sgc.modules {
		laststate := sgc.modules[i].present
		sgc.modules[i].present = sgc.modules[i].mctrl.TestIfPresent()
		if laststate != sgc.modules[i].present && sgc.modules[i].present {
			sgc.ModFullUpdate(i)
		}
	}
}
