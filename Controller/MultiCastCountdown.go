package controller

import (
	"errors"
	"fmt"
	"net"
)

type MultiCastCountdown struct {
	con net.Conn
}

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
func (m *MultiCastCountdown) SendTimeStrike(time uint32, numStrike uint8) error {
	_, err := m.con.Write([]byte(fmt.Sprintf("{timeleft:%d,strike:%d}", time, numStrike)))
	return err
}
func (m *MultiCastCountdown) SendWin() error {
	_, err := m.con.Write([]byte("{win:true}"))
	return err
}
func (m *MultiCastCountdown) SendBoom() error {
	_, err := m.con.Write([]byte("{boom:true}"))
	return err
}
func (m *MultiCastCountdown) SendReset() error {
	_, err := m.con.Write([]byte("{win:false,boom:false,timeleft:0,strike:0}"))
	return err
}
