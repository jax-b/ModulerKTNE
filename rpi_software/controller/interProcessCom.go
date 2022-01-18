package controller

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	ipc "github.com/james-barrow/golang-ipc"
	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"

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
	sv, err := ipc.StartServer("ktne-ipc", nil)
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

func (sipc *InterProcessCom) Run() {
	go func() {
		for {
			select {
			case quit := <-sipc.closech:
				if quit {
					return
				}
			default:
				sipc.commandTree()
			}
		}
	}()
}

// Safely closes the interprocess communicator
func (sipc *InterProcessCom) Close() error {
	sipc.closech <- true
	sipc.ipc.Close()
	return nil
}

func (sipc *InterProcessCom) commandTree() {
	if sipc.ipc.StatusCode() == 3 {
		message, err := sipc.ipc.Read()
		if err != nil {
			sipc.log.Error("Failed to read from IPC", err)
		}
		messages := string(message.Data)
		sipc.log.Debug("Received message from IPC", message)
		messagesCMDDTA := strings.Split(messages, ":")
		messagesCMD := strings.Split(messagesCMDDTA[0], ".")
		var ipcwerr error
		if message.MsgType == 3 {
			switch messagesCMD[1] {
			case "game":
				switch messagesCMD[2] {
				case "start":
					err := sipc.game.StartGame()
					if err == nil {
						sipc.log.Info("Game Started")
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.start.ok"))
					} else {
						sipc.log.Error("Failed to Start Game:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.start.error"))
					}
				case "stop":
					err := sipc.game.StopGame()
					if err == nil {
						sipc.log.Info("Game Stopped")
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.stop.ok"))
					} else {
						sipc.log.Error("Failed to Stop game:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.stop.error"))
					}
				case "set_time": // Attempts to set the time. This command will be automatically followed by get_time
					gametime, err := time.ParseDuration(messagesCMDDTA[1])
					if err != nil {
						sipc.log.Errorf("Failed to convert time: %e", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_time.error"))
						break
					}
					err = sipc.game.SetGameTime(gametime)
					if err == nil {
						sipc.log.Info("Set game time to:", gametime)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_time.ok"))
					} else {
						sipc.log.Error("Failed to set time:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_time.error"))
					}
					fallthrough
				case "get_time":
					gametime := sipc.game.GetGameTime()
					buffer := []byte("mktne.game.time:")
					buffer = append(buffer, []byte(strconv.Itoa(int(gametime)))...)
					ipcwerr = sipc.ipc.Write(2, buffer)
				case "set_strike": // Attempts to set the strike count. This command will be automatically followed by get_strike
					strike, err := strconv.ParseInt(messagesCMDDTA[1], 10, 16)
					if err != nil {
						sipc.log.Error("Failed to convert strike:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_strike.error"))
						break
					}
					err = sipc.game.SetStrikes(int8(strike))
					if err == nil {
						sipc.log.Info("Set strike to:", strike)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_strike.ok"))
					} else {
						sipc.log.Error("Failed to set strike:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_strike.error"))
					}
					fallthrough
				case "get_strike":
					strikes := sipc.game.GetStrikes()
					buffer := []byte("mktne.game.strike:")
					buffer = append(buffer, []byte(strconv.Itoa(int(strikes)))...)
					ipcwerr = sipc.ipc.Write(2, buffer)
				case "set_strike_rate": // Attempts to set the strike rate. This command will be automatically followed by get_strike_rate
					strikeRate, err := strconv.ParseFloat(messagesCMDDTA[1], 32)
					if err != nil {
						sipc.log.Error("Failed to convert strike rate:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_strike_rate.error"))
						break
					}
					err = sipc.game.SetStrikeRate(float32(strikeRate))
					if err == nil {
						sipc.log.Info("Set strike rate to:", strikeRate)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_strike_rate.ok"))
					} else {
						sipc.log.Error("Failed to set strike rate:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.set_strike_rate.error"))
					}
					fallthrough
				case "get_strike_rate":
					strikeRate := sipc.game.GetStrikeRate()
					buffer := []byte("mktne.game.strike_rate:")
					buffer = append(buffer, []byte(strconv.FormatFloat(float64(strikeRate), 'f', 2, 32))...)
					ipcwerr = sipc.ipc.Write(2, buffer)
				case "set_serialnumber": // Attempts to set the serial number. This command will be automatically followed by get_serialnumber
					serialNumber := messagesCMDDTA[1]
					err := sipc.game.SetSerial(serialNumber)
					if err != nil {
						sipc.log.Error("Failed to set the serial number", err)
						sipc.ipc.Write(1, []byte("mktne.game.set_serialnumber.error"))
						break
					}
					sipc.log.Info("Set serial number to:", serialNumber)
					sipc.ipc.Write(1, []byte("mktne.game.set_serialnumber.ok"))
					fallthrough
				case "get_serialnumber":
					sipc.ipc.Write(1, []byte("mktne.game.serialnumber:"+sipc.game.GetSerial()))
				case "add_indicator": // Attempts to add a indicator to the list. This command will be automatically followed by get_indicators.
					var indiobj Indicator
					err := json.Unmarshal([]byte(messagesCMDDTA[1]), &indiobj)
					if err != nil {
						sipc.log.Error("Failed to unmarshal indicator:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.add_indicator.error"))
						break
					}
					sipc.game.AddIndicator(indiobj)
					sipc.log.Info("Added Indicator to list", indiobj)
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.add_indicator.ok"))
					fallthrough
				case "get_indicators":
					indicators := sipc.game.GetIndicators()
					buffer := []byte("mktne.game.indicators:")
					indjson, err := json.Marshal(indicators)
					if err != nil {
						sipc.log.Error("Failed to marshal indicators json")
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.get_indicators.error"))
						break
					}
					buffer = append(buffer, indjson...)
					ipcwerr = sipc.ipc.Write(2, buffer)
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.get_indicators.ok"))
				case "clear_indicators":
					sipc.game.ClearIndicators()
					sipc.log.Info("Cleared Active Indicators")
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.clear_indicators.ok"))
				case "add_port": // Attempts to add a port to the list. This command will be automatically followed by get_ports.
					portInt64, err := strconv.ParseInt(messagesCMDDTA[1], 10, 8)
					if err != nil {
						sipc.log.Error("Failed to convert port:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.add_port.error"))
						break
					}
					port := byte(portInt64)
					err = sipc.game.SetPorts(port)
					if err != nil {
						sipc.log.Error("Failed to add port:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.add_port.error"))
						break
					}
					sipc.log.Info("Added Port to list", port)
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.add_port.ok"))
					fallthrough
				case "get_ports":
					ports := sipc.game.GetPorts()
					buffer := []byte("mktne.game.ports:")
					portjson, err := json.Marshal(ports)
					if err != nil {
						sipc.log.Error("Failed to marshal ports json")
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.get_ports.error"))
						break
					}
					buffer = append(buffer, portjson...)
					ipcwerr = sipc.ipc.Write(2, buffer)
				case "clear_ports":
					sipc.game.ClearPorts()
					sipc.log.Info("Cleared Game Ports")
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.game.clear_ports.ok"))
				}
			case "module":
				modnum64, err := strconv.ParseInt(messagesCMD[4], 10, 8)
				if err != nil {
					sipc.log.Error("Failed to convert module number:", err)
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.module.error"))
					break
				}
				modnum := int(modnum64)
				switch messagesCMD[3] {
				case "get_present":
					buffer := []byte("mktne.module." + messagesCMD[4] + ".present:")
					if sipc.game.modules[modnum].present {
						buffer = append(buffer, []byte("true")...)
					} else {
						buffer = append(buffer, []byte("false")...)
					}
					ipcwerr = sipc.ipc.Write(2, buffer)
				case "set_seed":
					seed, err := strconv.ParseInt(messagesCMDDTA[1], 10, 16)
					if err != nil {
						sipc.log.Error("Failed to convert seed:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".set_seed.error"))
						break
					}
					err = sipc.game.modules[modnum].mctrl.SetGameSeed(uint16(seed))
					if err == nil {
						sipc.log.Info("Set module ", modnum, " seed to: ", seed)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".set_seed.ok"))
					} else {
						sipc.log.Error("Failed to set module ", modnum, " seed: ", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.module."+messagesCMD[4]+".set_seed.error"))
					}
				case "get_type":
					buffer := []byte("mktne.module." + messagesCMD[4] + ".type:")
					mtype := sipc.game.modules[modnum].modtype
					var mtypesl []rune
					for i := range mtype {
						mtypesl = append(mtypesl, mtype[i])
					}
					buffer = append(buffer, []byte(string(mtypesl))...)
					ipcwerr = sipc.ipc.Write(2, buffer)
				}
			case "network":
				switch messagesCMD[2] {
				case "close":
					sipc.game.multicast.mnetc.Close()
					sipc.game.multicast.useMulti = false
					ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.close.ok"))
					sipc.log.Info("Closed multicast")
				case "open":
					sipc.game.multicast.useMulti = true
					sipc.game.multicast.mnetc, err = mktnecf.NewMultiCastCountdown(sipc.game.log, sipc.game.cfg.Network.MultiCastIP, sipc.game.cfg.Network.MultiCastPort)
					if err != nil {
						sipc.log.Error("Failed to open multicast:", err)
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.open.error"))
					} else {
						sipc.log.Info("Opened multicast")
						ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.open.ok"))
					}
				case "change":
					switch messagesCMD[3] {
					case "port":
						port, err := strconv.ParseInt(messagesCMDDTA[1], 10, 32)
						if err != nil {
							sipc.log.Error("Failed to convert port:", err)
							ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.change.port.error"))
							break
						}
						err = sipc.game.multicast.mnetc.ChangePort(int(port))
						if err != nil {
							sipc.log.Error("Failed to change port:", err)
							ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.change.port.error"))
						} else {
							sipc.log.Info("Changed port to:", port)
							ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.change.port.ok"))
						}
					case "ip":
						err := sipc.game.multicast.mnetc.ChangeIP(messagesCMDDTA[1])
						if err != nil {
							sipc.log.Error("Failed to change IP:", err)
							ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.change.ip.error"))
						} else {
							sipc.log.Info("Changed MCast IP to:", messagesCMDDTA[1])
							ipcwerr = sipc.ipc.Write(1, []byte("mktne.network.change.ip.ok"))
						}
					}
				}
			}
		}
		if ipcwerr != nil {
			sipc.log.Fatal("Failed to write to IPC:", ipcwerr)
		}
	}
}

func (sipc *InterProcessCom) SyncStatus(stat *mktnecf.Status) error {
	type msg struct {
		Time                string  `json:"timeleft"`
		NumStrike           int8    `json:"strike"`
		Boom                bool    `json:"boom"`
		Win                 bool    `json:"win"`
		Gamerun             bool    `json:"gamerun"`
		Strikereductionrate float32 `json:"strikerate"`
	}
	omsg := msg{
		Time:                stat.Time.String(),
		NumStrike:           int8(stat.NumStrike),
		Boom:                stat.Boom,
		Win:                 stat.Win,
		Gamerun:             stat.Gamerun,
		Strikereductionrate: stat.Strikereductionrate,
	}
	json, err := json.Marshal(omsg)
	if err != nil {
		sipc.log.Warn("Failed to marshal status: ", err)
		return err
	}
	ipcwerr := sipc.ipc.Write(9, json)
	if ipcwerr != nil {
		sipc.log.Fatal("Failed to write to IPC:", ipcwerr)
		return ipcwerr
	}
	return nil
}
