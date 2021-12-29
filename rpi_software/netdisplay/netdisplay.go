package netdisplay

import (
	"math"
	"os"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
	"go.uber.org/zap"
)

var (
	tlastupdate time.Duration
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
		Time:                time.Hour,
		NumStrike:           0,
		Boom:                false,
		Win:                 false,
		Gamerun:             false,
		Strikereductionrate: 0,
	}

	tstop := make(chan bool)

	nd := &netdisplay{
		cstatus: cstatus,
		cscreen: "home",
		net:     listener,
		log:     log,
		tstop:   tstop,
		UI:      ui,
	}
	// Force a refresh of the UI just in case
	tlastupdate = 5 * time.Hour
	// Create the game timer
	tmr := mktnecf.NewGameTimer(log, nd.cstatus)
	// Set the callback function
	tmr.AddCallbackFunction(nd.timerCallbackFunc)

	// Store the timer in the netdisplay
	nd.tmr = tmr

	return nd
}

func (snd *netdisplay) Run() {
	// Start UI, Net OPS, Timer, and Network Packet Tree
	snd.UI.StartUI()
	snd.net.Run()
	go snd.tmr.Run()
	go snd.incomingPacketTree()
}

// Stops the timer, Net OPS, and UI
func (snd *netdisplay) Close() {
	go func() { snd.tmr.Close() }()
	snd.net.Close()
	snd.UI.Close()
}

// Processes incoming packets
func (snd *netdisplay) incomingPacketTree() {
	// Subscribe to the incoming packets
	statchan := snd.net.Subscribe()
	for {
		// Read a pakcet from the channel
		status := <-statchan
		// Update the status
		snd.cstatus.Time = status.Time
		snd.cstatus.Boom = status.Boom
		snd.cstatus.Gamerun = status.Gamerun
		snd.cstatus.Win = status.Win
		// If a end connition is met set the ui to that condition
		if snd.cstatus.Boom {
			snd.cscreen = "boom"
		} else if snd.cstatus.Win {
			snd.cscreen = "win"
		} else {
			// If the  game is running set the ui to that condition else go to the home screen
			if snd.cstatus.Gamerun {
				snd.cscreen = "gametime"
			} else {
				snd.cscreen = "home"
			}
		}
		// Update the Rates and Strikes
		snd.cstatus.Strikereductionrate = status.Strikereductionrate
		snd.cstatus.NumStrike = status.NumStrike
		// Force a refresh of the UI for the next timer cycle
		tlastupdate = 5 * time.Hour
		// Update the UI
		newmsg, _ := snd.UI.createMSG(snd.cstatus.Time, snd.cscreen, snd.cstatus.NumStrike)
		snd.UI.UpdateUI(newmsg)
	}
}

func (snd *netdisplay) timerCallbackFunc(tsb time.Time, tsa time.Time, stat mktnecf.Status) {
	snd.log.Debug("TimeLeft:", stat.Time)
	// If we still have a minute left
	if stat.Time.Minutes() > 1 {
		// Check to see if we have a different number of seconds
		if (tlastupdate - stat.Time).Seconds() >= 1 {
			// Send the update to the UI
			newmsg, _ := snd.UI.createMSG(stat.Time, snd.cscreen, stat.NumStrike)
			snd.log.Info("Sending new message:", string(newmsg[:]))
			snd.UI.UpdateUI(newmsg)
			// Update the last update time
			tlastupdate = stat.Time
		}
	} else {
		// If we don't have a minute left
		// Check to see if we have a different number of hundredths of a second to display
		if math.Mod((tlastupdate-stat.Time).Seconds(), 0.1) >= 0.01 {
			// Send the update to the UI
			newmsg, _ := snd.UI.createMSG(stat.Time, snd.cscreen, stat.NumStrike)
			snd.log.Info("Sending new message:", string(newmsg[:]))
			snd.UI.UpdateUI(newmsg)
			// Update the last update time
			tlastupdate = stat.Time
		}
	}
	// Check for tmr runnout
	if stat.Boom {
		snd.log.Info("Boom!")
		snd.cscreen = "boom"
		snd.cstatus.Boom = true
		snd.cstatus.Gamerun = false
		// Update UI
		newmsg, _ := snd.UI.createMSG(stat.Time, snd.cscreen, stat.NumStrike)
		snd.UI.UpdateUI(newmsg)
	}
}
