package netdisplay

import (
	"os"
	"time"

	"go.uber.org/zap"
)

type netdisplay struct {
	cstatus *Status
	cscreen string
	net     *MultiCastListener
	log     *zap.SugaredLogger
	UI      *UI
	tstop   chan bool
	lastmsg string
}

func NewNetDisplay(log *zap.SugaredLogger) *netdisplay {
	// Create and load the config
	appconfig := NewConfig(log)
	appconfig.Load()
	log.Info("Config loaded")
	// Create Listener and exit if failed
	listener, err := NewMultiCastListener(log, appconfig)
	if err != nil {
		log.Error("Failed to create listener", zap.Error(err))
		os.Exit(1)
	}

	// Create UI
	ui := NewUI(log)

	// Create the current status tracker
	cstatus := Status{
		Time:                time.Hour,
		NumStrike:           0,
		Boom:                false,
		Win:                 false,
		Gamerun:             false,
		Strikereductionrate: 0,
	}

	tstop := make(chan bool)
	return &netdisplay{
		cstatus: &cstatus,
		cscreen: "home",
		net:     listener,
		log:     log,
		tstop:   tstop,
		UI:      ui,
	}
}

func (snd *netdisplay) Run() {
	snd.UI.StartUI()
	snd.net.Run()
	go GameTimer(snd)
	go snd.incomingPacketTree()
}

func (snd *netdisplay) Close() {
	snd.net.Close()
	go func() { snd.tstop <- true }()
	snd.UI.Close()
}

func (snd *netdisplay) incomingPacketTree() {
	statchan := snd.net.Subscribe()
	for {
		status := <-statchan
		if snd.cstatus.Time != status.Time {
			snd.cstatus.Time = status.Time
		}
		if snd.cstatus.Boom != status.Boom {
			snd.cstatus.Boom = status.Boom
			if snd.cstatus.Boom {
				snd.cscreen = "boom"
			} else {
				if snd.cstatus.Gamerun {
					snd.cscreen = "gametime"
				} else {
					snd.cscreen = "home"
				}
			}
		}
		if snd.cstatus.Gamerun != status.Gamerun {
			snd.cstatus.Gamerun = status.Gamerun
			if snd.cstatus.Gamerun {
				snd.cscreen = "gametime"
			} else {
				snd.cscreen = "home"
			}
		}
		if snd.cstatus.Win != status.Win {
			snd.cstatus.Win = status.Win
			snd.cscreen = "win"
			if snd.cstatus.Win {
				snd.cscreen = "win"
			} else {
				if snd.cstatus.Gamerun {
					snd.cscreen = "gametime"
				} else {
					snd.cscreen = "home"
				}
			}
		}
		if snd.cstatus.Strikereductionrate != status.Strikereductionrate {
			snd.cstatus.Strikereductionrate = status.Strikereductionrate
		}
		if snd.cstatus.NumStrike != status.NumStrike {
			snd.cstatus.NumStrike = status.NumStrike
		}
		newmsg, _ := snd.UI.createMSG(snd.cstatus.Time, snd.cscreen, snd.cstatus.NumStrike)
		snd.UI.UpdateUI(newmsg)
	}
}
