package controller

import (
	"errors"
	"math/rand"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"

	"go.uber.org/zap"
)

var (
	VALID_INDICATORS []string = []string{"SND", "CLR", "CAR", "IND", "FRQ", "SIG", "NSA", "MSA", "TRN", "BOB", "FRK"}
	VALID_PORT_ID    []string = []string{"DVI", "PAR", "PS2", "RJ4", "SER", "RCA"}
)

const (
	GAMEPLAYMAXNUMPORT    = 6
	GAMEPLAYMAXTINDICATOR = 12
)

type module struct {
	mctrl   *mktnecf.ModControl
	present bool
	solved  bool
	modtype [4]rune
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
	timer      *mktnecf.GameTimer
	indicators []Indicator
	port       byte
	serialnum  [8]rune
	numbat     int
	maxstrike  uint8
}
type GameController struct {
	sidePanel      *SideControl
	modules        [10]module
	multicast      multicast
	game           gameinfo
	rpishield      *ShieldControl
	cfg            *Config
	btnWatchStopCh chan bool
	interStopCh    chan bool
	solvedStopCh   chan bool
	log            *zap.SugaredLogger
	rnd            *rand.Rand

	audio struct {
		closech chan bool
		tick    chan bool
	}
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
			maxstrike: 3,
		},
	}

	// Create the Side Panel control object

	gc.sidePanel = NewSideControl(gc.log, int(gc.cfg.Shield.I2cBusNumber))

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

	// Setup Audio Timer Tick
	gc.audio.closech = make(chan bool)
	gc.audio.tick = make(chan bool)

	gc.rpishield.TimerBeep(gc.audio.closech, gc.audio.tick)

	// Write Idle to the screen
	gc.rpishield.WriteIdle()

	// Set up timer
	gc.game.timer = mktnecf.NewGameTimer(gc.log, &gc.game.comStat)
	gc.game.timer.AddCallbackFunction(gc.tmrCallbackFunction)

	// Initalize RNG
	var src mktnecf.CryptoSource
	rnd := rand.New(src)
	gc.rnd = rnd

	return gc
}

// Starts Game Controller Monitoring components
func (sgc *GameController) Run() {
	sgc.rpishield.Run()
	go sgc.game.timer.Run()
	sgc.buttonWatcher()
	sgc.m2cInterruptHandler()
	sgc.solvedCheck()
}

// Safe Shutdown of all components
func (sgc *GameController) Close() {
	sgc.log.Info("Closing Game Controller")
	// Stopping interupt Watchers
	go func() { sgc.btnWatchStopCh <- true }()
	go func() { sgc.interStopCh <- true }()
	// Stopping Solved Watcher
	go func() { sgc.solvedStopCh <- true }()
	// Stopping Audio Ticker
	go func() { sgc.audio.closech <- true }()
	// Stopping Timer
	sgc.game.timer.Close()
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
	sgc.sidePanel.Close()

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
}

// Checks to see if each module is solved and if all of them are solved then trigger the win condition
func (sgc *GameController) solvedCheck() {
	for {
		select {
		case <-sgc.solvedStopCh:
			return
		default:
			solved := true
			for i := range sgc.modules {
				if sgc.modules[i].present && !sgc.modules[i].solved {
					solved = false
				}
			}
			if solved && sgc.game.comStat.Gamerun {
				sgc.log.Info("Game Has Been Won!")
				sgc.GameOverWin()
				if sgc.multicast.useMulti {
					sgc.multicast.mnetc.SendStatus(&sgc.game.comStat)
				}
			}
		}
	}
}

// Starts the game
func (sgc *GameController) StartGame() error {
	// for all the modules that are present, start the game
	sgc.scanAllModules()
	noModPresent := true
	for i := range sgc.modules {
		if sgc.modules[i].present {
			noModPresent = false
			sgc.modules[i].solved = false
			sgc.modules[i].mctrl.StartGame()
		}
	}
	if !noModPresent {
		sgc.game.comStat.Gamerun = true
		return nil
	} else {
		return errors.New("noModulesPresent")
	}

}

// Stops the game
func (sgc *GameController) StopGame() error {
	sgc.game.comStat.Gamerun = false
	// for all the modules that are present, stop the game
	for i := range sgc.modules {
		if sgc.modules[i].present {
			sgc.modules[i].mctrl.StopGame()
		}
	}

	return nil
}

// populates each module with a random game
func (sgc *GameController) RandomPopulate() {
	// Serial number generation
	serialLen := sgc.rnd.Intn(8-6) + 6
	charset := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := range sgc.game.serialnum {
		if i < serialLen {
			charnum := sgc.rnd.Intn(len(charset))
			sgc.game.serialnum[i] = charset[charnum]
		} else {
			sgc.game.serialnum[i] = 0
		}
	}

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
		sgc.game.indicators = append(sgc.game.indicators, indi)
	}

	// Port Generation
	sgc.game.port = byte(sgc.rnd.Intn(63))

	sgc.scanAllModules()
	for i := range sgc.modules {
		err := sgc.ModFullUpdate(i)
		if err != nil {
			sgc.log.Error("Could not update module: ", err)
		}
	}
}
