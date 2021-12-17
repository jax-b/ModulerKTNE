package netdisplay

import "time"

// this function is the timekeeper for the game
// it will stop the game when the time runs out
func Timer(StopCh chan bool, cstatus *Status) {
	ticker := time.NewTicker(time.Millisecond)
	extratick := 0
	for {
		if cstatus.Gamerun {
			select {
			case <-StopCh:
				return
			case <-ticker.C:
				// Need to add reduction rate
				cstatus.Time--
				if cstatus.NumStrike < 0 {
					everyrate := int((1 / cstatus.Strikereductionrate) / (-1 * float32(cstatus.Strikereductionrate)))
					if extratick >= everyrate {
						if cstatus.Time > 0 {
							cstatus.Time--
							extratick = 0
						} else {
							cstatus.Boom = true
							cstatus.Gamerun = false
						}
					} else {
						extratick++
					}
				}
			}
		}
	}
}
