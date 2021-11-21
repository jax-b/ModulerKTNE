package controller

import (
	"github.com/jax-b/ModulerKTNE/controller"
)

type GameController struct {
	sidePanels [4]*Controller.SideControl
	modules    [10]struct {
		mctrl   *Controller.ModControl
		present bool
	}
	rpishield *Controller.ShieldControl
	ipc       *Controller.InterProcessCom
	mnetc     *Controller.MultiCastCountdown
}
