package netdisplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	mktnecf "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"
)

type MultiCastListener struct {
	con         *ipv4.PacketConn
	log         *zap.SugaredLogger
	conf        *Config
	exitLoop    chan bool
	subscribers []chan mktnecf.Status
	interfaces  []net.Interface
	listpac     net.PacketConn
}

// Creates a new multicast connection
func NewMultiCastListener(logger *zap.SugaredLogger, cfg *Config) (*MultiCastListener, error) {
	logger = logger.Named("MulticastUtil")

	// Initalize the mcast listener struct
	mcastl := &MultiCastListener{
		log:  logger,
		conf: cfg,
	}

	// Check if the IP is valid
	if !net.ParseIP(cfg.Network.MultiCastIP).IsMulticast() {
		return nil, errors.New("IP is not a multicast address: " + cfg.Network.MultiCastIP)
	}

	// Create a new packet listener
	var err error
	mcastl.listpac, err = net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%d", cfg.Network.MultiCastPort))
	if err != nil {
		logger.Error("Failed to create Packet Listener:", err)
		return nil, err
	}
	// Create a new Connection
	mcastl.con = ipv4.NewPacketConn(mcastl.listpac)

	// Grab all of the interfaces
	mcastl.interfaces, err = net.Interfaces()
	if err != nil {
		logger.Error("Failed to get interfaces:", err)
		os.Exit(1)
	}

	// Set up multicast on each one
	for _, infa := range mcastl.interfaces {
		if infa.Flags&net.FlagUp == net.FlagUp && infa.Flags&net.FlagMulticast == net.FlagMulticast {
			logger.Debug("Interface is up and multicast capable: ", infa.Name)
			logger.Info("Joining interface to IGMP Group: ", infa.Name)
			if err = mcastl.con.JoinGroup(&infa, &net.UDPAddr{IP: net.ParseIP(cfg.Network.MultiCastIP)}); err != nil {
				logger.Error("Failed to Join IGMP Group on interface:", infa.Name, cfg.Network.MultiCastIP, err)
				return nil, err
			}
		}
	}
	logger.Info("Listening for status messages on ", cfg.Network.MultiCastIP, ":", cfg.Network.MultiCastPort)
	return mcastl, nil
}

// Creates a new Status subscriber and returns the status chanel
func (smcl *MultiCastListener) Subscribe() chan mktnecf.Status {
	newchan := make(chan mktnecf.Status)
	smcl.subscribers = append(smcl.subscribers, newchan)
	return newchan
}

// Starts up the loop for processing status messages
func (smcl *MultiCastListener) Run() {
	smcl.exitLoop = make(chan bool)
	go func() {
		for {
			select {
			case <-smcl.exitLoop:
				close(smcl.exitLoop)
				return
			default:
				smcl.log.Info("Waiting for status")
				status, srcaddr, err := smcl.getStatus()

				if err != nil {
					if err.(net.Error).Error() == "use of closed network connection" {
						return
					}
					smcl.log.Error("Failed to Process Status Message: ", err)
				} else if status != nil {
					smcl.log.Info("Got status from "+srcaddr.String()+": ", status)
					for _, subscriber := range smcl.subscribers {
						go func(sub chan mktnecf.Status) {
							sub <- *status
						}(subscriber)
					}
				}
			}
		}
	}()
}

func (smcl *MultiCastListener) Close() {
	go func() {
		smcl.exitLoop <- true
	}()
	// Tear down IGMP multicast on each interface
	for _, infa := range smcl.interfaces {
		if infa.Flags == (net.FlagMulticast & net.FlagUp) {
			if err := smcl.con.LeaveGroup(&infa, &net.UDPAddr{IP: net.ParseIP(smcl.conf.Network.MultiCastIP)}); err != nil {
				smcl.log.Error("Failed to leave IGMP group on interaface", infa.Name, smcl.conf.Network.MultiCastIP, err)
			}
		}
	}
	smcl.con.Close()
	smcl.listpac.Close()
}

// Gets current status as a json string from the multicast address
func (smcl *MultiCastListener) getStatus() (*mktnecf.Status, net.Addr, error) {
	buffer := make([]byte, 100)
	numread, _, src, err := smcl.con.ReadFrom(buffer)
	if err != nil {
		return nil, nil, err
	}

	// Make Sure we get enough data
	if numread > 5 {

		// Log the buffer contents
		smcl.log.Debug("Buffer Contents: ", string(buffer[:]))

		// Unmarshal the json data from the buffer
		status := &mktnecf.Status{}
		err = json.Unmarshal(buffer[:numread], status)
		if err != nil {
			smcl.log.Error("Error unmarshalling json: " + err.Error())
			return nil, nil, err
		}
		status.Time, _ = time.ParseDuration(status.StrTime)
		// Return the status message
		return status, src, nil
	} else {
		return nil, nil, errors.New("Failed to enough data from the network")
	}
}
