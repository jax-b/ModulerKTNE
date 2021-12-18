package netdisplay

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// this function is the timekeeper for the game
// it will stop the game when the time runs out
func GameTimer(nd *netdisplay) {
	log := nd.log.Named("Time")
	nd.cscreen = "gametime"
	tlast := time.Now()
	tick := time.NewTicker(time.Millisecond)
	for {
		// Wait for the ticker to tick -- This is just to make sure we don't go to fast
		<-tick.C
		// if the timer is not running exit this iteration
		if !nd.cstatus.Gamerun {
			break
		}

		// Store wht seconds before subtracting
		tsb := nd.cstatus.Time.Seconds()
		// Store the current time
		tc := time.Now()
		// Subtract the time
		nd.cstatus.Time -= tc.Sub(tlast)
		// If we have stikes calculate the number of extra time to remove per strike and remove it
		if nd.cstatus.NumStrike > 0 {
			everyrate := int((1 / 0.25) / float32(nd.cstatus.NumStrike))
			nd.cstatus.Time -= tc.Sub(tlast) / time.Duration(everyrate)
		}
		// Save the time that we did all this math
		tlast = tc
		// Store the seconds after subtracting
		tsa := nd.cstatus.Time.Seconds()
		go func() {
			// If we still have a minute left
			if nd.cstatus.Time.Minutes() > 1 {
				// Check to see if we have a different number of seconds
				if int(tsb) != int(tsa) {
					// Send the update to the UI
					newmsg, _ := nd.UI.createMSG(nd.cstatus.Time, nd.cscreen, nd.cstatus.NumStrike)
					log.Debug("Sending new message:", string(newmsg[:]))
					nd.UI.UpdateUI(newmsg)
				}
			} else {
				// If we don't have a minute left
				// Check to see if we have a different number of hundreths of a second to display
				if fmt.Sprintf("%0.2f", tsb) != fmt.Sprintf("%0.2f", tsa) {
					// Send the update to the UI
					newmsg, _ := nd.UI.createMSG(nd.cstatus.Time, nd.cscreen, nd.cstatus.NumStrike)
					log.Debug("Sending new message:", string(newmsg[:]))
					nd.UI.UpdateUI(newmsg)
				}
			}
		}()
		go timeChecker(nd, log)

		// If we have been called to exit then exit
		select {
		case <-nd.tstop:
			close(nd.tstop)
			return
		default:
		}
	}
}
func timeChecker(nd *netdisplay, log *zap.SugaredLogger) {
	if nd.cstatus.Time <= 0 || nd.cstatus.Time.Minutes() > 100 {
		nd.cstatus.Boom = true
		nd.cstatus.Gamerun = false
		nd.cscreen = "boom"
		newmsg, _ := nd.UI.createMSG(time.Duration(0), nd.cscreen, nd.cstatus.NumStrike)
		log.Info("Sending new message:", string(newmsg[:]))
		nd.UI.UpdateUI(newmsg)
	}
}
