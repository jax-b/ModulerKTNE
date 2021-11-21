package controller

import (
	"github.com/james-barrow/golang-ipc"

	"github.com/jax-b/ModulerKTNE/rpi_software/controller"
	"go.uber.org/zap"
)

type InterProcessCom struct {
	ipc     *ipc.Server
	closech chan bool
	game    *Controller.GameController
	logger  *zap.Logger
}

// Creates a new interprocess communicator
func NewIPC(gamectrl *Controller.GameController) (*InterProcessCom, error) {
	scon := &ipc.ServerConfig{
		Timeout: time.Millisecond * 50,
	}
	sv, err := ipc.StartServer("ktne-ipc", scon)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ipc := &InterProcessCom{
		ipc:     sv,
		closech: make(chan bool),
		game:    gamectrl,
	}
	return ipc, nil
}

// Safely closes the interprocess communicator
func (self *InterProcessCom) Close() error {
	self.closech <- true
	self.ipc.Close()
	return nil
}

func (self *InterProcessCom) commandTree() {
	if self.ipc.StatusCode() = 3{
		message, err := self.ipc.Read()
		if err != nil {
			log.Println(err)
		}
		message = strings.Split(message)
		switch message[0] {
			case "mktne.start_game":
				err := self.game.StartGame()
				if err == nil {
					self.ipc.Write([]byte("mktne.start_game.ok"))
				} else {
					self.logger.error("Failed to Start Game:", err)
					self.ipc.Write([]byte("mktne.start_game.error"))
				}
				break
			case "mktne.stop_game":
				err := self.game.StopGame()
				if err == nil {
					self.ipc.Write([]byte("mktne.stop_game.ok"))
				} else {
					self.logger.error("Failed to Stop ipc:", err)
					self.ipc.Write([]byte("mktne.stop_game.error"))
				}
				break
			case "mktne.get_time":
				gametime := self.game.get_time()
				buffer := []byte("mktne.gametime:")
				buffer = buffer.append(buffer,[]byte(strconv.Itoa(gametime)))
				self.ipc.Write(buffer)
				break
			case "mktne.set_time":
				gametime, err := strconv.ParseInt(message[1],10,16)
				if err != nil {
					self.logger.error("Failed to convert time:", err)
					self.ipc.Write([]byte("mktne.set_time.error"))
				}
				err = self.game.set_time(gametime)
				if err == nil {
					self.ipc.Write([]byte("mktne.set_time.ok"))
				} else {
					self.logger.error("Failed to set time:", err)
					self.ipc.Write([]byte("mktne.set_time.error"))
				}
				break
		}
	}
}

func (self *InterProcessCom) timerRunOut() {
	self.ipc.Write([]byte("mktne.timer_runout"))
}