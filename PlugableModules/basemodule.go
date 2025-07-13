package basemodule

import (
	"fmt"
	"machine"

	"encoding/base64"

	"github.com/jax-b/ModulerKTNE/WireTypes"
	"google.golang.org/protobuf/proto"
)

const (
	SuccessLEDPin     = machine.GP8
	AddressInPin      = machine.GP26
	RandSourcePin     = machine.GP27
	FailureLEDPin     = machine.GP9
	S2MInterruptPin   = machine.GP4
	ControllerAddress = 0x00
	BroadcastAddress  = 0xFF
)

type ModuleController struct {
	GameState     *WireTypes.GameState
	ModuleState   *WireTypes.ModuleState
	ModuleAddress uint32
}

func NewModuleController() *ModuleController {
	return &ModuleController{
		GameState: &WireTypes.GameState{},
	}
}

func (mc *ModuleController) SendModuleInfo(
	destination uint32,
	operation WireTypes.Operation,
) {

	message := &WireTypes.WireMessage{
		Sender:   mc.ModuleAddress,
		Receiver: destination,
		Option:   operation,
	}
	message.Content = &WireTypes.WireMessage_ModuleStateContent{
		ModuleStateContent: mc.ModuleState,
	}
	data, err := proto.Marshal(message)

	if err != nil {
		fmt.Printf("Error marshalling message: %v\n", err)
		return
	}

	fmt.Print(base64.StdEncoding.EncodeToString(data))
}
