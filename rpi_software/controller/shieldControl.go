package controller

import (
	"fmt"
	"time"

	"github.com/jax-b/go-i2c7Seg"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
)

type ShieldControl struct {
	strikecount         uint8
	seg                 *i2c7Seg.SevenSegI2C
	buzzerPin           rpio.Pin
	strike1Pin          rpio.Pin
	strike2Pin          rpio.Pin
	modInterruptPin     rpio.Pin
	mfbPin              rpio.Pin
	mfbCallbackConsumer []chan uint16
	m2cCallbackConsumer []chan bool
	mfbEdge             rpio.Edge
	stopBtnCheck        chan bool
	log                 *zap.SugaredLogger
}

func NewShieldControl(logger *zap.SugaredLogger, cfg *Config) *ShieldControl {
	logger = logger.Named("ShieldControl")
	i2c, err := i2c7Seg.NewSevenSegI2C(cfg.Shield.SevenSegAddress, int(cfg.Shield.I2cBusNumber))
	if err != nil {
		logger.Error("Failed to create i2c7Seg", err)
	}
	err = rpio.Open()
	if err != nil {
		logger.Error("Failed to open rpio", err)
	}
	sc := &ShieldControl{
		seg:             i2c,
		strikecount:     0,
		buzzerPin:       rpio.Pin(cfg.Shield.BuzzerPinNum),
		strike1Pin:      rpio.Pin(cfg.Shield.Strike1PinNum),
		strike2Pin:      rpio.Pin(cfg.Shield.Strike2PinNum),
		modInterruptPin: rpio.Pin(cfg.Shield.ModInterruptPinNum),
		mfbPin:          rpio.Pin(cfg.Shield.MfbStartPinNum),
		log:             logger,
	}
	// Configure pins
	sc.buzzerPin.Mode(rpio.Pwm)
	sc.strike1Pin.Output()
	sc.strike2Pin.Output()
	sc.mfbPin.Input()
	sc.mfbPin.PullUp()
	sc.mfbPin.Detect(rpio.RiseEdge)
	sc.mfbEdge = rpio.RiseEdge
	sc.modInterruptPin.Input()
	sc.modInterruptPin.PullUp()
	sc.modInterruptPin.Detect(rpio.FallEdge)

	// Start Input checker
	go func() {
		stop := false
		for !stop {
			sc.btnCheck()
			select {
			case stop = <-sc.stopBtnCheck:
			default:
			}
		}
	}()
	return sc
}

// Closes out all functions that are running safely
func (self *ShieldControl) Close() {
	self.stopBtnCheck <- true
	self.ClearDisplay()
	self.seg.Close()
	rpio.Close()
}

// Adds a strike to the display system and plays a sound
func (self *ShieldControl) AddStrike() {
	self.strikecount++
	// span strike sound in a seprate concurent
	go self.buzzStrikeSound()
	if self.strikecount == 1 {
		self.strike1Pin.High()
		self.strike2Pin.Low()
	} else if self.strikecount == 2 {
		self.strike1Pin.High()
		self.strike2Pin.High()
	} else {
		self.strike1Pin.Low()
		self.strike2Pin.High()
	}
}

// Resets the strike to zero
func (self *ShieldControl) ResetStrike() {
	self.strikecount = 0
	self.strike1Pin.Low()
	self.strike2Pin.Low()
}

// plays the sound of a strike
func (self *ShieldControl) buzzStrikeSound() {
	self.buzzerPin.Freq(64000)
	self.buzzerPin.DutyCycle(16, 32)
	time.Sleep(time.Millisecond * 500)
	self.buzzerPin.DutyCycle(0, 32)
}

// TimeClockSignal
func (self *ShieldControl) TimeSigBeep(multiplier float32) {
	self.buzzerPin.Freq(600)
	self.buzzerPin.DutyCycle(16, 32)
	time.Sleep(time.Millisecond * time.Duration(350*multiplier))
	self.buzzerPin.Freq(60)
	time.Sleep(time.Millisecond * time.Duration(250*multiplier))
	self.buzzerPin.DutyCycle(0, 32)
}

// Writes the name of the game to the display
func (self *ShieldControl) WriteIdle() {
	self.seg.Clear()
	self.seg.WriteAsciiChar(0, 'K', false)
	self.seg.WriteAsciiChar(1, 'T', true)
	self.seg.WriteAsciiChar(3, 'N', false)
	self.seg.WriteAsciiChar(4, 'E', true)
	self.seg.WriteDisplay()
}

// Converts and writes the time to the display max it can display is 99:60
// Will move to 49.50 when time is less then 1 minute
func (self *ShieldControl) WriteTime(timemilis uint32) {
	self.seg.Clear()
	var tstring string
	timemilisf := float32(timemilis)
	MinutesRemaining := int(timemilisf*0.001) / 60
	SecondsRemaining := int(timemilisf*0.001) % 60
	if MinutesRemaining > 0 {
		if MinutesRemaining > 99 {
			MinutesRemaining = 99
		}
		tstring = fmt.Sprintf("%2d:%02d", MinutesRemaining, SecondsRemaining)
		trune := []rune(tstring)
		for i := uint8(0); i < 5; i++ {
			self.seg.WriteAsciiChar(i, byte(trune[i]), false)
		}
		self.seg.DrawColon(true)
	} else {
		hundtensec := int(timemilisf*0.1) % 100
		tstring = fmt.Sprintf("%2d.%02d", SecondsRemaining, hundtensec)

		trune := []rune(tstring)
		self.seg.WriteAsciiChar(0, byte(trune[0]), false)
		self.seg.WriteAsciiChar(1, byte(trune[1]), true)
		self.seg.WriteAsciiChar(3, byte(trune[3]), false)
		self.seg.WriteAsciiChar(4, byte(trune[4]), false)
	}
	self.seg.WriteDisplay()
}

// Clears the display
func (self *ShieldControl) ClearDisplay() {
	self.seg.Clear()
	self.seg.WriteDisplay()
}

// Registers a consumer of the module to controller interrupt line
// The consumer will have a true sent down it when the interupt line is triggered
func (self *ShieldControl) RegisterM2CConsumer() chan bool {
	c := make(chan bool)
	self.m2cCallbackConsumer = append(self.m2cCallbackConsumer, c)

	return c
}

// Registers a consumer of the module to the multifunction button
// Sends how long the button was held for
func (self *ShieldControl) RegisterMFBConsumer() chan uint16 {
	c := make(chan uint16)
	self.mfbCallbackConsumer = append(self.mfbCallbackConsumer, c)

	return c
}

// Function for checking if the external inputs were pressed and will signal consumers when ready
func (self *ShieldControl) btnCheck() {
	// if we have a signal from a downstream controller signal all consumers (nonblocking)
	if self.modInterruptPin.EdgeDetected() {
		for _, c := range self.m2cCallbackConsumer {
			go func(c chan bool) {
				c <- true
			}(c)
		}
	}
	// if the button state is changed, and we are looking for a press
	if self.mfbPin.EdgeDetected() && self.mfbEdge == rpio.RiseEdge {
		// Span a new concurent to wait for release
		go func() {
			// Set up the edge detection
			self.mfbEdge = rpio.FallEdge
			self.mfbPin.Detect(rpio.FallEdge)
			// Record time of press
			mfbPush := time.Now()
			// Wait for release
			for !self.mfbPin.EdgeDetected() {
			}
			// Record time of release and compute difference
			mfbRelease := time.Now()
			mfbPushTime := uint16(mfbRelease.Sub(mfbPush).Milliseconds())
			// Signal all consumers for how long the button was pressed
			for _, c := range self.mfbCallbackConsumer {
				// Non blocking channel update
				go func(c chan uint16) {
					c <- mfbPushTime
				}(c)
			}
			// Reset edge detection
			self.mfbEdge = rpio.RiseEdge
			self.mfbPin.Detect(rpio.RiseEdge)
		}()
	}
}
