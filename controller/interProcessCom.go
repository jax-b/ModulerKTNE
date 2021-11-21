package controller

import (
	"rpc"

	"github.com/jax-b/ModulerKTNE/Controller"
)

type InterProcessCom struct {
	rpc   *rpc.RPC
	close chan bool
	game  *Controller.GameController
}

func NewIPC(gamectrl *Controller.GameController) *InterProcessCom {
	return &InterProcessCom{
		rpc:   rpc.NewRPC(),
		close: make(chan bool),
		game:  gamectrl,
	}
}

func (ipc *InterProcessCom) Close() error {
	ipc.close <- true
	return nil
}
