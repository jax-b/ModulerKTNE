package netdisplay

import "time"

// this function is the timekeeper for the game
// it will stop the game when the time runs out
func GameTimer(StopCh chan bool, nd *netdisplay) {
	ticker := time.NewTicker(time.Millisecond)
	extratick := 0
	for {
		if nd.cstatus.Gamerun {
			select {
			case <-StopCh:
				return
			case <-ticker.C:
				// Need to add reduction rate
				nd.cstatus.Time--
				if nd.cstatus.NumStrike < 0 {
					everyrate := int((1 / nd.cstatus.Strikereductionrate) / (-1 * float32(nd.cstatus.Strikereductionrate)))
					if extratick >= everyrate {
						if nd.cstatus.Time > 0 {
							nd.cstatus.Time--
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
			if nd.cstatus.Time <= 0 {
				nd.cstatus.Boom = true
				nd.cstatus.Gamerun = false
				nd.cscreen = "boom"
			}
		}
		newmsg, _ := nd.UI.createMSG(timetostring(nd.cstatus.Time), nd.cscreen, nd.cstatus.NumStrike)
		if nd.lastmsg != string(newmsg[:]) {
			nd.UI.UpdateUI(newmsg)
			nd.lastmsg = string(newmsg[:])
		}
	}
}
