package commonfiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type MultiCastCountdown struct {
	con net.Conn
	log *zap.SugaredLogger
}

// Creates a new multicast connection
func NewMultiCastCountdown(logger *zap.SugaredLogger, MultiCastIP string, MulticastPort int) (*MultiCastCountdown, error) {
	var newcon net.Conn
	var err error
	logger = logger.Named("MulticastUtil")
	ip := net.ParseIP(MultiCastIP)
	udpaddr, _ := net.ResolveUDPAddr("udp", ip.String()+":"+fmt.Sprint(MulticastPort))
	if ip.IsMulticast() {
		newcon, err = net.DialUDP("udp", nil, udpaddr)
		logger.Info("Multicasting on: ", udpaddr)
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
func (smcc *MultiCastCountdown) SendStatus(stat *Status) error {
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
		return err
	}
	_, err = smcc.con.Write(json)
	return err
}

// Resets the segments to zero via a json string to the multicast address
func (smcc *MultiCastCountdown) SendReset() error {
	err := smcc.SendStatus(&Status{
		Time:                time.Duration(0),
		NumStrike:           0,
		Boom:                false,
		Win:                 false,
		Gamerun:             false,
		Strikereductionrate: 0.25,
	})
	return err
}
