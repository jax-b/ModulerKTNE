package controller

func (sgc *GameController) updateNetworkIPC() error {
	// err := sgc.ipc.SyncStatus(&sgc.game.comStat)
	// if err != nil {
	// 	return err
	// }
	var err error
	if sgc.multicast.useMulti {
		err = sgc.multicast.mnetc.SendStatus(&sgc.game.comStat)
	}
	return err
}
