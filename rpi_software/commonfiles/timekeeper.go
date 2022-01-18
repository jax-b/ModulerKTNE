package commonfiles

import (
	"time"

	"go.uber.org/zap"
)

type GameTimer struct {
	stat        *Status
	log         *zap.SugaredLogger
	callbacksub []func(tsb time.Time, tsa time.Time, stat Status)
	tstop       chan bool
}

func NewGameTimer(logger *zap.SugaredLogger, stat *Status) *GameTimer {
	log := logger.Named("Time")

	return &GameTimer{
		stat:  stat,
		log:   log,
		tstop: make(chan bool),
	}
}

// this function is the timekeeper for the game
// it will stop the game when the time runs out
func (sgt *GameTimer) Run() {
	tlast := time.Now()
	tick := time.NewTicker(time.Millisecond)
	for {
		// If we have been called to exit then exit
		select {
		case <-sgt.tstop:
			sgt.log.Info("TimeKeeper Stopped")
			close(sgt.tstop)
			return
		case <-tick.C: // Wait for the ticker to tick -- This is just to make sure we don't go to fast
			// if the timer is not running exit this iteration
			if !sgt.stat.Gamerun {
				continue
			}
			// sgt.log.Debug("Time: ", sgt.stat.Time)

			// Store the current time
			tstart := time.Now()
			// Subtract the time
			// if tstart.Sub(tlast) > time.Second {
			// 	continue
			// }
			sgt.stat.Time -= tstart.Sub(tlast)

			// If we have stikes calculate the number of extra time to remove per strike and remove it
			if sgt.stat.NumStrike > 0 {
				everyrate := (1 / sgt.stat.Strikereductionrate) / float32(sgt.stat.NumStrike)
				var textra time.Duration
				if everyrate < 1 {
					textra = tstart.Sub(tlast)
					textra += time.Duration(float32(textra.Nanoseconds()) * (1 - everyrate))
				} else {
					textra = tstart.Sub(tlast) / time.Duration(everyrate)
				}
				sgt.log.Infof("Rate: %.2f,Time Extra: %s", everyrate, textra.String())
				sgt.stat.Time -= textra
			}
			// Save the time that we did all this math
			tlast = tstart

			// Send the time to the callback functions
			for _, callback := range sgt.callbacksub {
				go callback(tstart, time.Now(), *sgt.stat)
			}

			// Check if the game is over
			go sgt.timeChecker()
		}
	}
}

func (sgt *GameTimer) Close() {
	sgt.tstop <- true
}

func (sgt *GameTimer) AddCallbackFunction(callback func(tsb time.Time, tsa time.Time, stat Status)) {
	sgt.log.Debug("Adding callback function")
	sgt.callbacksub = append(sgt.callbacksub, callback)
}

func (sgt *GameTimer) timeChecker() {
	if sgt.stat.Time <= 0 || sgt.stat.Time.Minutes() > 100 {
		sgt.stat.Boom = true
		sgt.stat.Gamerun = false
		sgt.log.Debug("Game over: Boom Protocol Activated")
		for _, callback := range sgt.callbacksub {
			go callback(time.Now(), time.Now(), *sgt.stat)
		}
	}
}
