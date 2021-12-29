package netdisplay

import (
	"math"
	"os"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
	"go.uber.org/zap"
)

var (
	tlastupdate time.Time
)

type netdisplay struct {
	cstatus *mktnecf.Status
	cscreen string
	net     *MultiCastListener
	log     *zap.SugaredLogger
	UI      *UI
	tstop   chan bool
	lastmsg string
	tmr     *mktnecf.GameTimer
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
	cstatus := &mktnecf.Status{
		Time:                time.Minute,
		NumStrike:           0,
		Boom:                false,
		Win:                 false,
		Gamerun:             true,
		Strikereductionrate: 0,
	}

	tstop := make(chan bool)

	nd := &netdisplay{
		cstatus: cstatus,
		cscreen: "gametime",
		net:     listener,
		log:     log,
		tstop:   tstop,
		UI:      ui,
	}
	tmr := mktnecf.NewGameTimer(log, nd.cstatus)
	nd.tmr = tmr

	tmr.AddCallbackFunction(nd.timerCallbackFunc)

	return nd
}

func (snd *netdisplay) Run() {
	snd.UI.StartUI()
	snd.net.Run()
	go snd.tmr.Run()
	go snd.incomingPacketTree()
	tlastupdate = time.Now()
	msg, _ := snd.UI.createMSG(snd.cstatus.Time, snd.cscreen, snd.cstatus.NumStrike)
	snd.UI.UpdateUI(msg)
}

func (snd *netdisplay) Close() {
	go func() { snd.tmr.Close() }()
	snd.net.Close()
	snd.UI.Close()
}

func (snd *netdisplay) incomingPacketTree() {
	statchan := snd.net.Subscribe()
	for {
		status := <-statchan
		snd.cstatus.Time = status.Time
		snd.cstatus.Boom = status.Boom
		snd.cstatus.Gamerun = status.Gamerun
		snd.cstatus.Win = status.Win
		if snd.cstatus.Boom {
			snd.cscreen = "boom"
		} else if snd.cstatus.Win {
			snd.cscreen = "win"
		} else {
			if snd.cstatus.Gamerun {
				snd.cscreen = "gametime"
			} else {
				snd.cscreen = "home"
			}
		}
		snd.cstatus.Strikereductionrate = status.Strikereductionrate
		snd.cstatus.NumStrike = status.NumStrike
		newmsg, _ := snd.UI.createMSG(snd.cstatus.Time, snd.cscreen, snd.cstatus.NumStrike)
		snd.UI.UpdateUI(newmsg)
	}
}

func (snd *netdisplay) timerCallbackFunc(tsb time.Time, tsa time.Time, stat mktnecf.Status) {
	snd.log.Debug("TimeLeft:", stat.Time)
	// If we still have a minute left
	if stat.Time.Minutes() > 1 {

		// Check to see if we have a different number of seconds
		if tsa.Sub(tlastupdate).Seconds() >= 1 {
			// Send the update to the UI
			newmsg, _ := snd.UI.createMSG(stat.Time, snd.cscreen, stat.NumStrike)
			snd.log.Debug("Sending new message:", string(newmsg[:]))
			snd.UI.UpdateUI(newmsg)
			tlastupdate = tsa
		}
	} else {
		// If we don't have a minute left
		// Check to see if we have a different number of hundredths of a second to display
		if math.Mod(tlastupdate.Sub(tsa).Seconds(), 0.1)*-1 >= 0.01 {
			// Send the update to the UI
			newmsg, _ := snd.UI.createMSG(stat.Time, snd.cscreen, stat.NumStrike)
			snd.log.Debug("Sending new message:", string(newmsg[:]))
			snd.UI.UpdateUI(newmsg)
			tlastupdate = tsa
		}
	}
	if stat.Boom {
		snd.log.Info("Boom!")
		snd.cscreen = "boom"
		snd.cstatus.Boom = true
		snd.cstatus.Gamerun = false
		newmsg, _ := snd.UI.createMSG(stat.Time, snd.cscreen, stat.NumStrike)
		snd.UI.UpdateUI(newmsg)
	}
}
