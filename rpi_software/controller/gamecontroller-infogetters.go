package controller

import "time"

// Gets the current game time
func (sgc *GameController) GetGameTime() time.Duration {
	return sgc.game.comStat.Time
}

// Gets the current amount of strikes
func (sgc *GameController) GetStrikes() int8 {
	if sgc.game.comStat.Win {
		return 1
	}
	return int8(sgc.game.comStat.NumStrike) * -1
}

// Get the srike reduction rate
func (sgc *GameController) GetStrikeRate() float32 {
	return sgc.game.comStat.Strikereductionrate
}

// Gets the currently configured indicators
func (sgc *GameController) GetIndicators() []Indicator {
	return sgc.game.indicators
}

// returns all of the ports that are configured for the game
func (sgc *GameController) GetPorts() []byte {
	return sgc.game.port
}

func (sgc *GameController) GetSerial() string {
	return string(sgc.game.serialnum[0:])
}
