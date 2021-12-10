package testnetdata

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type MultiCastCountdown struct {
	con net.Conn
	log *zap.SugaredLogger
}

// Creates a new multicast connection
func NewMultiCastCountdown(logger *zap.SugaredLogger, cfg *Config) (*MultiCastCountdown, error) {
	var newcon net.Conn
	var err error
	logger = logger.Named("MulticastUtil")
	ip := net.ParseIP(cfg.Network.MultiCastIP)
	if ip.IsMulticast() {
		newcon, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: ip, Port: cfg.Network.MultiCastPort})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("IP is not a multicast address: " + ip.String())
	}
	return &MultiCastCountdown{con: newcon, log: logger}, nil
}

func (smcc *MultiCastCountdown) Close() {
	smcc.con.Close()
}

func (smcc *MultiCastCountdown) ChangeIP(ip string) error {
	nip := net.ParseIP(ip)
	if !nip.IsMulticast() {
		return errors.New("IP is not a multicast address: " + ip)
	}
	port := strings.Split(smcc.con.LocalAddr().String(), ":")[1]
	portint, err := strconv.ParseInt(port, 10, 16)
	if err != nil {
		return err
	}
	smcc.con.Close()
	smcc.con, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: nip, Port: int(portint)})
	return err
}

func (smcc *MultiCastCountdown) ChangePort(port int) error {
	ip := strings.Split(smcc.con.LocalAddr().String(), ":")[0]
	smcc.con.Close()
	var err error
	smcc.con, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IP(ip), Port: port})
	return err
}

// Sends current status as a json string to the multicast address
func (smcc *MultiCastCountdown) SendStatus(time uint32, numStrike int8, boom bool, win bool) error {
	wins := "false"
	booms := "false"
	if win {
		wins = "true"
	}
	if boom {
		booms = "true"
	}
	_, err := smcc.con.Write([]byte(fmt.Sprintf("{timeleft:%d,strike:%d,win:%s,boom:%s}", time, numStrike, wins, booms)))
	return err
}

// Resets the segments to zero via a json string to the multicast address
func (smcc *MultiCastCountdown) SendReset() error {
	_, err := smcc.con.Write([]byte("{win:false,boom:false,timeleft:0,strike:0}"))
	return err
}
