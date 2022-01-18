package controller

// Adds a indicator to the list
func (sgc *GameController) AddIndicator(indi Indicator) {
	if len(sgc.game.indicators) > GAMEPLAYMAXTINDICATOR {
		sgc.game.indicators[len(sgc.game.indicators)] = indi
	} else {
		sgc.game.indicators = append(sgc.game.indicators, indi)
	}
	if indi.Lit {
		for i := range sgc.modules {
			if sgc.modules[i].present {
				sgc.modules[i].mctrl.SetGameLitIndicator(indi.Label)
			}
		}
	}
}

// Clears out the current indicators
func (sgc *GameController) ClearIndicators() {
	for i := range sgc.modules {
		if sgc.modules[i].present {
			sgc.modules[i].mctrl.ClearGameLitIndicator()
		}
	}
	sgc.game.indicators = make([]Indicator, 0)
}

// clears all of the ports that are configured for the game
func (sgc *GameController) ClearPorts() error {
	sgc.game.port = 0x0
	for i := range sgc.modules {
		if sgc.modules[i].present {
			err := sgc.modules[i].mctrl.ClearGamePortIDS()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
