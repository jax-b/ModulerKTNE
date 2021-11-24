package controller

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	ipc "github.com/james-barrow/golang-ipc"

	"go.uber.org/zap"
)

type InterProcessCom struct {
	ipc     *ipc.Server
	closech chan bool
	game    *GameController
	log     *zap.SugaredLogger
}

// Creates a new interprocess communicator
func NewIPC(logger *zap.SugaredLogger, gamectrl *GameController) *InterProcessCom {
	logger = logger.Named("IPC")
	scon := &ipc.ServerConfig{
		Timeout: time.Millisecond * 50,
	}
	sv, err := ipc.StartServer("ktne-ipc", scon)
	if err != nil {
		log.Fatal("Failed to start IPC Server", err)
	}
	ipc := &InterProcessCom{
		ipc:     sv,
		closech: make(chan bool),
		game:    gamectrl,
		log:     logger,
	}
	return ipc
}

// Safely closes the interprocess communicator
func (self *InterProcessCom) Close() error {
	self.closech <- true
	self.ipc.Close()
	return nil
}

func (self *InterProcessCom) commandTree() {
	if self.ipc.StatusCode() == 3 {
		message, err := self.ipc.Read()
		if err != nil {
			self.log.Error("Failed to read from IPC", err)
		}
		messages := string(message.Data)
		self.log.Debug("Received message from IPC", message)
		messagesCMDDTA := strings.Split(messages, ":")
		messagesCMD := strings.Split(messagesCMDDTA[0], ".")
		if message.MsgType == 3 {
			switch messagesCMD[1] {
			case "game":
				switch messagesCMD[2] {
				case "start":
					err := self.game.StartGame()
					if err == nil {
						self.log.Info("Game Started")
						self.ipc.Write(1, []byte("mktne.game.start.ok"))
					} else {
						self.log.Error("Failed to Start Game:", err)
						self.ipc.Write(1, []byte("mktne.game.start.error"))
					}
					break
				case "stop":
					err := self.game.StopGame()
					if err == nil {
						self.log.Info("Game Stopped")
						self.ipc.Write(1, []byte("mktne.game.stop.ok"))
					} else {
						self.log.Error("Failed to Stop game:", err)
						self.ipc.Write(1, []byte("mktne.game.stop.error"))
					}
					break
				case "set_time": // Attempts to set the time. This command will be automatically followed by get_time
					gametime, err := strconv.ParseInt(messagesCMDDTA[1], 10, 32)
					if err != nil {
						self.log.Error("Failed to convert time:", err)
						self.ipc.Write(1, []byte("mktne.game.set_time.error"))
						break
					}
					err = self.game.SetGameTime(uint32(gametime))
					if err == nil {
						self.log.Info("Set game time to:", gametime)
						self.ipc.Write(1, []byte("mktne.game.set_time.ok"))
					} else {
						self.log.Error("Failed to set time:", err)
						self.ipc.Write(1, []byte("mktne.game.set_time.error"))
						break
					}
				case "get_time":
					gametime := self.game.GetGameTime()
					buffer := []byte("mktne.game.time:")
					buffer = append(buffer, []byte(strconv.Itoa(int(gametime)))...)
					self.ipc.Write(2, buffer)
					break
				case "set_strike": // Attempts to set the strike count. This command will be automatically followed by get_strike
					strike, err := strconv.ParseInt(messagesCMDDTA[1], 10, 16)
					if err != nil {
						self.log.Error("Failed to convert strike:", err)
						self.ipc.Write(1, []byte("mktne.game.set_strike.error"))
						break
					}
					err = self.game.SetStrikes(int8(strike))
					if err == nil {
						self.log.Info("Set strike to:", strike)
						self.ipc.Write(1, []byte("mktne.game.set_strike.ok"))
					} else {
						self.log.Error("Failed to set strike:", err)
						self.ipc.Write(1, []byte("mktne.game.set_strike.error"))
						break
					}
				case "get_strike":
					strikes := self.game.GetStrikes()
					buffer := []byte("mktne.game.strike:")
					buffer = append(buffer, []byte(strconv.Itoa(int(strikes)))...)
					self.ipc.Write(2, buffer)
					break
				case "set_strike_rate":
					strikeRate, err := strconv.ParseFloat(messagesCMDDTA[1], 32)
					if err != nil {
						self.log.Error("Failed to convert strike rate:", err)
						self.ipc.Write(1, []byte("mktne.game.set_strike_rate.error"))
						break
					}
					err = self.game.SetStrikeRate(float32(strikeRate))
					if err == nil {
						self.log.Info("Set strike rate to:", strikeRate)
						self.ipc.Write(1, []byte("mktne.game.set_strike_rate.ok"))
					} else {
						self.log.Error("Failed to set strike rate:", err)
						self.ipc.Write(1, []byte("mktne.game.set_strike_rate.error"))
						break
					}
				}
			case "module":
				modnum64, err := strconv.ParseInt(messagesCMD[4], 10, 8)
				if err != nil {
					self.log.Error("Failed to convert module number:", err)
					self.ipc.Write(1, []byte("mktne.module.error"))
					break
				}
				modnum := int(modnum64)
				switch messagesCMD[3] {
				case "get_present":
					buffer := []byte("mktne.module." + messagesCMD[4] + ".present:")
					if self.game.modules[modnum].present {
						buffer = append(buffer, []byte("true")...)
					} else {
						buffer = append(buffer, []byte("false")...)
					}
					self.ipc.Write(2, buffer)
				case "set_seed":
					seed, err := strconv.ParseInt(messagesCMDDTA[1], 10, 16)
					if err != nil {
						self.log.Error("Failed to convert seed:", err)
						self.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".set_seed.error"))
						break
					}
					err = self.game.modules[modnum].mctrl.SetGameSeed(uint16(seed))
					if err == nil {
						self.log.Info("Set module", modnum, "seed to:", seed)
						self.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".set_seed.ok"))
					} else {
						self.log.Error("Failed to set module", modnum, "seed:", err)
						self.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".set_seed.error"))
					}
					break
				case "get_type":
					buffer := []byte("mktne.module." + messagesCMD[4] + ".type:")
					mtype, err := self.game.modules[modnum].mctrl.GetModuleType()
					if err != nil {
						self.log.Error("Failed to get module type:", err)
						self.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".get_type.error"))
						break
					}
					var mtypesl []rune
					for i := range mtype {
						mtypesl = append(mtypesl, mtype[i])
					}
					buffer = append(buffer, []byte(string(mtypesl))...)
					self.ipc.Write(2, buffer)
					break
				}
			case "network":
				switch messagesCMD[2] {
				case "close":
					self.game.multicast.mnetc.Close()
					self.game.multicast.useMulti = false
					self.ipc.Write(1, []byte("mktne.network.close.ok"))
					self.log.Info("Closed multicast")
					break
				case "open":
					self.game.multicast.useMulti = true
					self.game.multicast.mnetc, err = NewMultiCastCountdown(self.game.log, self.game.cfg)
					if err != nil {
						self.log.Error("Failed to open multicast:", err)
						self.ipc.Write(1, []byte("mktne.network.open.error"))
					} else {
						self.log.Info("Opened multicast")
						self.ipc.Write(1, []byte("mktne.network.open.ok"))
					}
					break
				case "change":
					switch messagesCMD[3] {
					case "port":
						port, err := strconv.ParseInt(messagesCMDDTA[1], 10, 32)
						if err != nil {
							self.log.Error("Failed to convert port:", err)
							self.ipc.Write(1, []byte("mktne.network.change.port.error"))
							break
						}
						err = self.game.multicast.mnetc.ChangePort(int(port))
						if err != nil {
							self.log.Error("Failed to change port:", err)
							self.ipc.Write(1, []byte("mktne.network.change.port.error"))
						} else {
							self.log.Info("Changed port to:", port)
							self.ipc.Write(1, []byte("mktne.network.change.port.ok"))
						}
					case "ip":
						err := self.game.multicast.mnetc.ChangeIP(messagesCMDDTA[1])
						if err != nil {
							self.log.Error("Failed to change IP:", err)
							self.ipc.Write(1, []byte("mktne.network.change.ip.error"))
						} else {
							self.log.Info("Changed MCast IP to:", messagesCMDDTA[1])
							self.ipc.Write(1, []byte("mktne.network.change.ip.ok"))
						}
						break
					}
				}
			}
		}
	}
}

func (self *InterProcessCom) SyncStatus(time uint32, numStrike int8, boom bool, win bool) {
	wins := "false"
	booms := "false"
	if win {
		wins = "true"
	}
	if boom {
		booms = "true"
	}
	statusjson := fmt.Sprintf("mktne.status:{timeleft:%d,strike:%d,win:%s,boom:%s}", time, numStrike, wins, booms)
	self.ipc.Write(9, []byte(statusjson))
}
