package controller

import (
	"time"
)

// Set the number of strikes
func (sgc *GameController) SetStrikes(strikes int8) error {
	var err error
	if strikes > 0 {
		sgc.game.comStat.Win = true
		err = sgc.updateModSolvedStatus(true)
	} else {
		sgc.game.comStat.NumStrike = uint8(strikes * -1)
		err = sgc.updateModSolvedStatus(false)
	}
	if err != nil {
		return err
	}
	err = sgc.updateNetworkIPC()
	return err
}

// Adds a strike to the current number of strikes
func (sgc *GameController) AddStrike() error {
	var err error
	sgc.game.comStat.NumStrike++
	if sgc.game.comStat.NumStrike >= sgc.game.maxstrike {
		sgc.game.comStat.Boom = true
		sgc.game.comStat.Gamerun = false
		sgc.game.comStat.Win = false
		err = sgc.StopGame()
		if err != nil {
			return err
		}
	} else {
		err = sgc.updateModSolvedStatus(false)
		if err != nil {
			return err
		}
		sgc.rpishield.AddStrike()
	}
	err = sgc.updateNetworkIPC()
	return err
}

// Sets the game time to the given time
func (sgc *GameController) SetGameTime(time time.Duration) error {
	sgc.game.comStat.Time = time
	err := sgc.updateModTime()
	if err != nil {
		return err
	}
	err = sgc.updateNetworkIPC()
	return err
}

// Set the strike reduction rate
func (sgc *GameController) SetStrikeRate(rate float32) error {
	sgc.game.comStat.Strikereductionrate = rate
	err := sgc.updateModStrikeRate()
	if err != nil {
		return err
	}
	err = sgc.updateNetworkIPC()
	return err
}

// Sets the current game serial number
func (sgc *GameController) SetSerial(serial string) error {
	for i := range serial {
		if i > len(sgc.game.serialnum) {
			break
		}
		sgc.game.serialnum[i] = rune(serial[i])
	}
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetGameSerialNumber(sgc.game.serialnum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
