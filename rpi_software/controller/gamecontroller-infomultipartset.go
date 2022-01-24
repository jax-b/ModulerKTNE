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
	if sgc.sidePanel.active {
		sgc.sidePanel.controller.SetIndicator(indi.Lit, indi.Label)
	}
}
