package netdisplay

import (
	"time"
)

// this function is the timekeeper for the game
// it will stop the game when the time runs out
func GameTimer(nd *netdisplay) {
	go timeChecker(nd)
	extratick := 0
	for range time.Tick(time.Millisecond) {
		if nd.cstatus.Gamerun {
			nd.cstatus.Time = nd.cstatus.Time - time.Millisecond
			if nd.cstatus.NumStrike < 0 {
				everyrate := int((1 / nd.cstatus.Strikereductionrate) / (-1 * float32(nd.cstatus.Strikereductionrate)))
				if extratick >= everyrate {
					if nd.cstatus.Time > 0 {
						nd.cstatus.Time = nd.cstatus.Time - time.Millisecond
						extratick = 0
					} else {
						nd.cstatus.Boom = true
						nd.cstatus.Gamerun = false
						nd.cscreen = "boom"
					}
				} else {
					extratick++
				}
			}
		}
	}
}
func timeChecker(nd *netdisplay) {
	log := nd.log.Named("TimeChecker")
	for {
		if nd.cstatus.Time <= 0 {
			nd.cstatus.Boom = true
			nd.cstatus.Gamerun = false
			nd.cscreen = "boom"
		}
		newmsg, _ := nd.UI.createMSG(nd.cstatus.Time, nd.cscreen, nd.cstatus.NumStrike)
		if nd.lastmsg != string(newmsg[:]) {
			log.Info("Sending new message:", string(newmsg[:]))
			nd.UI.UpdateUI(newmsg)
			nd.lastmsg = string(newmsg[:])
		}
	}
}
