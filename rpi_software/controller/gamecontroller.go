package controller

import (
	"math/rand"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"

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
	mctrl   *mktnecf.ModControl
	present bool
	solved  bool
}
type multicast struct {
	useMulti bool
	mnetc    *mktnecf.MultiCastCountdown
}
type Indicator struct {
	Lit   bool    `json:"lit"`
	Label [3]rune `json:"label"`
}
type gameinfo struct {
	comStat    mktnecf.Status
	indicators []Indicator
	port       []uint8
	serialnum  [8]rune
	numbat     int
	maxstrike  uint8
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

// Creates a new game controller takes in a variable called runAsDamon Which is a bool
// runAsDamon controls whether or not the
func NewGameCtrlr(log *zap.SugaredLogger) *GameController {
	// Load the configuration
	cfg := NewConfig(log)
	cfg.Load()

	// Try to open the RPI Shield Control
	rpis := NewShieldControl(log, cfg)

	// Create the GameController Object
	gc := &GameController{
		rpishield:      rpis,
		cfg:            cfg,
		log:            log,
		btnWatchStopCh: make(chan bool),
		timerStopCh:    make(chan bool),
		interStopCh:    make(chan bool),
		solvedStopCh:   make(chan bool),
		game: gameinfo{
			comStat: mktnecf.Status{
				Time:                time.Hour,
				NumStrike:           0,
				Gamerun:             false,
				Boom:                false,
				Strikereductionrate: 0.25,
				Win:                 false,
			},
			maxstrike: 2,
		},
	}

	// Create the inter process communicator object
	gc.ipc = NewIPC(log, gc)

	//Loop through the side panels and create their control objects
	SPADDR := [4]byte{TOP_PANEL, RIGHT_PANEL, BOTTOM_PANEL, LEFT_PANEL}
	for i := range gc.sidePanels {
		gc.sidePanels[i] = NewSideControl(gc.log, SPADDR[i], int(gc.cfg.Shield.I2cBusNumber))
	}

	//Loop through the modules and create their control objects
	MCADDR := [10]byte{mktnecf.FRONT_MOD_1, mktnecf.FRONT_MOD_2, mktnecf.FRONT_MOD_3, mktnecf.FRONT_MOD_4, mktnecf.FRONT_MOD_5, mktnecf.BACK_MOD_1, mktnecf.BACK_MOD_2, mktnecf.BACK_MOD_3, mktnecf.BACK_MOD_4, mktnecf.BACK_MOD_5}
	for i := range gc.modules {
		gc.modules[i].mctrl = mktnecf.NewModControl(gc.log, MCADDR[i], int(gc.cfg.Shield.I2cBusNumber))
	}

	// Check if multicast is enabled then create its object
	if gc.cfg.Network.UseMulticast {
		gc.multicast.useMulti = true
		var err error
		gc.multicast.mnetc, err = mktnecf.NewMultiCastCountdown(gc.log, gc.cfg.Network.MultiCastIP, gc.cfg.Network.MultiCastPort)
		if err != nil {
			gc.log.Error("Failed to create multicast countdown, proceeding without multicast", err)
			gc.multicast.useMulti = false
		}
	} else {
		gc.multicast.useMulti = false
	}

	// Initalize RNG
	var src mktnecf.CryptoSource
	rnd := rand.New(src)
	gc.rnd = rnd

	return gc
}

// Starts Game Controller Monitoring components
func (sgc *GameController) Run() {
	sgc.ipc.Run()
	sgc.rpishield.Run()
	go sgc.buttonWatcher()
	go sgc.m2cInterruptHandler()
	sgc.solvedCheck()
}

// Safe Shutdown of all components
func (sgc *GameController) Close() {
	sgc.log.Info("Closing Game Controller")
	go func() { sgc.timerStopCh <- true }()
	sgc.btnWatchStopCh <- true
	sgc.interStopCh <- true
	sgc.solvedStopCh <- true
	// flush the logger
	sgc.log.Sync()
	// Close all of the modules
	for _, mod := range sgc.modules {
		if mod.present {
			mod.mctrl.ClearAllGameData()
		}
		mod.mctrl.Close()
	}
	// Close the RPI Shield
	for i := range sgc.sidePanels {
		sgc.sidePanels[i].Close()
	}
	// Close the shield
	sgc.rpishield.Close()

	// Close multicast if used
	if sgc.multicast.useMulti {
		sgc.multicast.mnetc.SendStatus(&mktnecf.Status{
			Time:                time.Hour,
			NumStrike:           0,
			Gamerun:             false,
			Boom:                false,
			Strikereductionrate: 0.25,
			Win:                 false,
		})
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
			notSolved := true
			for i := range sgc.modules {
				if sgc.modules[i].present && !sgc.modules[i].solved {
					notSolved = false
				}
			}
			if !notSolved {
				sgc.game.comStat.Win = true
				sgc.game.comStat.Gamerun = false
				sgc.game.comStat.Boom = false
				sgc.StopGame()
				sgc.ipc.SyncStatus(&sgc.game.comStat)
				if sgc.multicast.useMulti {
					sgc.multicast.mnetc.SendStatus(&sgc.game.comStat)
				}
			}
		}
	}
}

// Handles the interrupt from the a modules updating its status in the game controller and updating strikes
func (sgc *GameController) m2cInterruptHandler() {
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
					if int16(solvedStat) < (int16(sgc.game.comStat.NumStrike) * -1) {
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

// Starts the game
func (sgc *GameController) StartGame() error {
	// for all the modules that are present, start the game
	sgc.scanAllModules()
	for i := range sgc.modules {
		if sgc.modules[i].present {
			sgc.modules[i].solved = false
			sgc.modules[i].mctrl.StartGame()
		}
	}
	sgc.game.comStat.Gamerun = true
	return nil
}

// Stops the game
func (sgc *GameController) StopGame() error {
	if sgc.game.comStat.Gamerun {
		// for all the modules that are present, stop the game
		for i := range sgc.modules {
			if sgc.modules[i].present {
				sgc.modules[i].mctrl.StopGame()
			}
		}
	}
	return nil
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
			Label: indilblrn,
			Lit:   indilit,
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
