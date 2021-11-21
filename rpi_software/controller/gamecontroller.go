package controller

import (
	"github.com/jax-b/ModulerKTNE/controller"
)

type GameController struct {
	sidePanels [4]*controller.SideControl
	modules    [10]struct {
		mctrl   *controller.ModControl
		present bool
	}
	rpishield *controller.ShieldControl
	ipc       *controller.InterProcessCom
	mnetc     *controller.MultiCastCountdown
	cfg       *controller.Config
}
