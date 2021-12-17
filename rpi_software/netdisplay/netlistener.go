package netdisplay

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
)

type MultiCastListener struct {
	con         net.Conn
	log         *zap.SugaredLogger
	exitLoop    chan bool
	subscribers []chan Status
}
type Status struct {
	Time                uint32  `json:"timeleft"`
	NumStrike           int8    `json:"strike"`
	Boom                bool    `json:"boom"`
	Win                 bool    `json:"win"`
	Gamerun             bool    `json:"gamerun"`
	Strikereductionrate float32 `json:"strikerate"`
}

// Creates a new multicast connection
func NewMultiCastListener(logger *zap.SugaredLogger, cfg *Config) (*MultiCastListener, error) {
	var newcon net.Conn
	var err error
	logger = logger.Named("MulticastUtil")
	ip := net.ParseIP(cfg.Network.MultiCastIP)
	if ip.IsMulticast() {
		newcon, err = net.ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: ip, Port: cfg.Network.MultiCastPort})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("IP is not a multicast address: " + ip.String())
	}
	return &MultiCastListener{con: newcon, log: logger}, nil
}

func (smcl *MultiCastListener) Subscribe() chan Status {
	newchan := make(chan Status)
	smcl.subscribers = append(smcl.subscribers, newchan)
	return newchan
}

func (smcl *MultiCastListener) Run() {
	smcl.exitLoop = make(chan bool)
	go func() {
		for {
			select {
			case <-smcl.exitLoop:
			default:
				status, err := smcl.getStatus()
				if err != nil {
					smcl.log.Error("Error getting status: ", err)
				} else {
					smcl.log.Debug("Got status: ", status)
					for _, subscriber := range smcl.subscribers {
						subscriber <- *status
					}
				}
			}
		}
	}()
}

func (smcl *MultiCastListener) Close() {
	smcl.con.Close()
}

// Sends current status as a json string to the multicast address
func (smcl *MultiCastListener) getStatus() (*Status, error) {
	buffer := make([]byte, 80)
	smcl.con.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	numread, err := smcl.con.Read(buffer)
	if err != nil {
		if err == os.ErrDeadlineExceeded {
			return nil, err
		}
		smcl.log.Error("Error reading from multicast address: " + err.Error())
		return nil, err
	} else {
		status := new(Status)
		err = json.Unmarshal(buffer[:numread], status)
		if err != nil {
			smcl.log.Error("Error unmarshalling json: " + err.Error())
			return nil, err
		}
		return status, nil
	}
}

// Resets the segments to zero via a json string to the multicast address
func (smcl *MultiCastListener) SendReset() error {
	_, err := smcl.con.Write([]byte("{win:false,boom:false,timeleft:0,strike:0}"))
	return err
}
