package controller

import (
	"errors"
	"fmt"
	"net"
)

type MultiCastCountdown struct {
	con net.Conn
}

// Creates a new multicast connection
func NewMultiCastCountdown(ip net.IP, port int) (*MultiCastCountdown, error) {
	var newcon net.Conn
	var err error
	if ip.IsMulticast() {
		newcon, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: ip, Port: port})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("IP is not a multicast address")
	}
	return &MultiCastCountdown{con: newcon}, nil
}

// Sends current status as a json string to the multicast address
func (m *MultiCastCountdown) SendTimeStrike(time uint32, numStrike uint8, boom bool, win, bool) error {
	if win {
		wins := "true"
	} else {}
		wins := "false"
	}
	if boom {
		booms := "true"
	} else {
		booms := "false"
	}
	_, err := m.con.Write([]byte(fmt.Sprintf("{timeleft:%d,strike:%d,win:%s,boom:%s}", time, numStrike, wins, booms)))
	return err
}
// Resets the segments to zero via a json string to the multicast address
func (m *MultiCastCountdown) SendReset() error {
	_, err := m.con.Write([]byte("{win:false,boom:false,timeleft:0,strike:0}"))
	return err
}
