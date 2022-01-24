package controller

func (sgc *GameController) ClearSerialNumber() {
	for i := range sgc.game.serialnum {
		sgc.game.serialnum[i] = ' '
	}
	for i := range sgc.modules {
		sgc.modules[i].mctrl.ClearGameSerialNumber()
	}
	if sgc.sidePanel.active {
		sgc.sidePanel.controller.ClearSerialNumber()
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
