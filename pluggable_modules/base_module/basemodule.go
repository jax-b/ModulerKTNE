package basemodule

import (
	"machine"
)

const (
	SuccessLEDPin   = machine.GP8
	AddressInPin    = machine.GP26
	RandSourcePin   = machine.GP27
	FailureLEDPin   = machine.GP9
	S2MInterruptPin = machine.GP4
)

type ModuleController struct {
	solved *GameState
}

func NewModuleController() *ModuleController {
	return &ModuleController{
		solved: &GameState,
	}
}
