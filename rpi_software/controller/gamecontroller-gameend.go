package controller

func (sgc *GameController) GameOverBoom() {
	sgc.game.comStat.Boom = true
	sgc.game.comStat.Gamerun = false
	sgc.game.comStat.Win = false
	sgc.StopGame()
	sgc.rpishield.ExploadSound()
	sgc.SetStrikes(0, true)
	sgc.ClearSeeds()
	sgc.rpishield.WriteIdle()
	for i := range sgc.modules {
		sgc.modules[i].solved = false
	}
}

func (sgc *GameController) GameOverWin() {
	sgc.game.comStat.Boom = false
	sgc.game.comStat.Gamerun = false
	sgc.game.comStat.Win = true
	sgc.StopGame()
	for i := range sgc.modules {
		sgc.modules[i].solved = false
	}
	sgc.rpishield.GameWinSound()
	sgc.SetStrikes(0, true)
	sgc.ClearSeeds()
	sgc.rpishield.WriteIdle()
}
