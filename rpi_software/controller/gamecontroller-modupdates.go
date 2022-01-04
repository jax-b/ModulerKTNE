package controller

import "errors"

// Polls all possible module addresses and sees if something is their. Updates the class variables
func (sgc *GameController) scanAllModules() {
	for i := range sgc.modules {
		laststate := sgc.modules[i].present
		sgc.modules[i].present = sgc.modules[i].mctrl.TestIfPresent()
		if laststate != sgc.modules[i].present && sgc.modules[i].present {
			sgc.ModFullUpdate(i)
		}
	}
}

// Updates all the modules to the stored number of strikes if the module is not solved
// Can be forced to update the modules even if the module is solved
func (sgc *GameController) updateModSolvedStatus(force bool) error {
	for _, mod := range sgc.modules {
		if mod.present && (!mod.solved || force) {
			var err error
			if sgc.game.comStat.Win {
				err = mod.mctrl.SetSolvedStatus(1)
			} else {
				err = mod.mctrl.SetSolvedStatus(int8(sgc.game.comStat.NumStrike) * -1)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Updates the time on all module to the current game time
func (sgc *GameController) updateModTime() error {
	for _, mod := range sgc.modules {
		if mod.present {
			err := mod.mctrl.SyncGameTime(uint32(sgc.game.comStat.Time.Milliseconds()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Updates the strike rate on all module to the current strike rate
func (sgc *GameController) updateModStrikeRate() error {
	for _, mod := range sgc.modules {
		if mod.present {
			err := mod.mctrl.SetStrikeReductionRate(sgc.game.comStat.Strikereductionrate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Updates the serial number on all modules to the current serial number
func (sgc *GameController) updateModSerial() error {
	for _, mod := range sgc.modules {
		if mod.present {
			err := mod.mctrl.SetGameSerialNumber(sgc.game.serialnum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Update the specified module with all of the game values
func (sgc *GameController) ModFullUpdate(modnum int) error {
	mod := sgc.modules[modnum]
	// Test if the module is present
	mod.present = mod.mctrl.TestIfPresent()
	if mod.present {
		var litindi [][3]rune
		for _, lbl := range sgc.game.indicators {
			if lbl.Lit {
				litindi = append(litindi, lbl.Label)
			}
		}
		err := mod.mctrl.SetupAllGameData(
			sgc.game.serialnum,
			litindi,
			uint8(sgc.game.numbat),
			sgc.game.port,
		)
		if err != nil {
			sgc.log.Error("Failed to update module's game data", err)
			return err
		}
		err = mod.mctrl.SetStrikeReductionRate(sgc.game.comStat.Strikereductionrate)
		if err != nil {
			sgc.log.Error("Failed to set strike reduction rate", err)
			return err
		}
		if sgc.game.comStat.Win {
			err = mod.mctrl.SetSolvedStatus(1)
		} else {
			err = mod.mctrl.SetSolvedStatus(int8(sgc.game.comStat.NumStrike) * -1)
		}
		return err
	} else {
		sgc.log.Errorf("Module %d not present", modnum)
		return errors.New("Module not present")
	}
}
