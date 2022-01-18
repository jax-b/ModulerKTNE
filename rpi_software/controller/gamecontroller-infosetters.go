package controller

import (
	"time"
)

// Set the number of strikes
func (sgc *GameController) SetStrikes(strikes int8, force ...bool) error {
	var forced bool = false
	if len(force) > 0 {
		forced = force[0]
	}

	var err error

	if strikes > 0 {
		sgc.game.comStat.Win = true
		err = sgc.updateModSolvedStatus(true)
	} else {
		sgc.game.comStat.NumStrike = uint8(strikes * -1)
		err = sgc.updateModSolvedStatus(forced)
	}

	if err != nil {
		return err
	}
	sgc.rpishield.SetStrike(uint8(strikes))
	err = sgc.updateNetworkIPC()
	return err
}

// Adds a strike to the current number of strikes
func (sgc *GameController) AddStrike() error {
	var err error
	sgc.game.comStat.NumStrike++
	if sgc.game.comStat.NumStrike >= sgc.game.maxstrike {
		sgc.log.Info("Number of strikes exceeded, BOOM!")
		sgc.GameOverBoom()
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
	sgc.updateModSerial()
	return nil
}

// Adds a port to the list
func (sgc *GameController) SetPorts(port byte) error {
	sgc.game.port = port
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetGamePortID(port)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sgc *GameController) SetNumBatteries(numbat uint8) error {
	sgc.game.numbat = int(numbat)
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.SetGameNumBatteries(numbat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sgc *GameController) SetModSeed(index uint8, seed uint16) {
	if index > 9 {
		return
	}
	sgc.modules[index].mctrl.SetGameSeed(seed)
}

func (sgc *GameController) SetMaxStrikes(inmax uint8) {
	sgc.game.maxstrike = inmax
}
