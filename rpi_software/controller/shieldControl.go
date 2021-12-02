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

	return sc
}

func (ssc *ShieldControl) Run() {
	// Start Input checker
	go func() {
		stop := false
		for !stop {
			ssc.btnCheck()
			select {
			case stop = <-ssc.stopBtnCheck:
			default:
			}
		}
	}()
}

// Closes out all functions that are running safely
func (ssc *ShieldControl) Close() {
	ssc.stopBtnCheck <- true
	ssc.ClearDisplay()
	ssc.seg.Close()
	rpio.Close()
}

// Adds a strike to the display system and plays a sound
func (ssc *ShieldControl) AddStrike() {
	ssc.strikecount++
	// span strike sound in a seprate concurent
	go ssc.buzzStrikeSound()
	if ssc.strikecount == 1 {
		ssc.strike1Pin.High()
		ssc.strike2Pin.Low()
	} else if ssc.strikecount == 2 {
		ssc.strike1Pin.High()
		ssc.strike2Pin.High()
	} else {
		ssc.strike1Pin.Low()
		ssc.strike2Pin.High()
	}
}

// Resets the strike to zero
func (ssc *ShieldControl) ResetStrike() {
	ssc.strikecount = 0
	ssc.strike1Pin.Low()
	ssc.strike2Pin.Low()
}

// plays the sound of a strike
func (ssc *ShieldControl) buzzStrikeSound() {
	ssc.buzzerPin.Freq(64000)
	ssc.buzzerPin.DutyCycle(16, 32)
	time.Sleep(time.Millisecond * 500)
	ssc.buzzerPin.DutyCycle(0, 32)
}

// TimeClockSignal
func (ssc *ShieldControl) TimeSigBeep(multiplier float32) {
	ssc.buzzerPin.Freq(600)
	ssc.buzzerPin.DutyCycle(16, 32)
	time.Sleep(time.Millisecond * time.Duration(350*multiplier))
	ssc.buzzerPin.Freq(60)
	time.Sleep(time.Millisecond * time.Duration(250*multiplier))
	ssc.buzzerPin.DutyCycle(0, 32)
}

// Writes the name of the game to the display
func (ssc *ShieldControl) WriteIdle() {
	ssc.seg.Clear()
	ssc.seg.WriteAsciiChar(0, 'K', false)
	ssc.seg.WriteAsciiChar(1, 'T', true)
	ssc.seg.WriteAsciiChar(3, 'N', false)
	ssc.seg.WriteAsciiChar(4, 'E', true)
	ssc.seg.WriteDisplay()
}

// Converts and writes the time to the display max it can display is 99:60
// Will move to 49.50 when time is less then 1 minute
func (ssc *ShieldControl) WriteTime(timemilis uint32) {
	ssc.seg.Clear()
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
			ssc.seg.WriteAsciiChar(i, byte(trune[i]), false)
		}
		ssc.seg.DrawColon(true)
	} else {
		hundtensec := int(timemilisf*0.1) % 100
		tstring = fmt.Sprintf("%2d.%02d", SecondsRemaining, hundtensec)

		trune := []rune(tstring)
		ssc.seg.WriteAsciiChar(0, byte(trune[0]), false)
		ssc.seg.WriteAsciiChar(1, byte(trune[1]), true)
		ssc.seg.WriteAsciiChar(3, byte(trune[3]), false)
		ssc.seg.WriteAsciiChar(4, byte(trune[4]), false)
	}
	ssc.seg.WriteDisplay()
}

// Clears the display
func (ssc *ShieldControl) ClearDisplay() {
	ssc.seg.Clear()
	ssc.seg.WriteDisplay()
}

// Registers a consumer of the module to controller interrupt line
// The consumer will have a true sent down it when the interupt line is triggered
func (ssc *ShieldControl) RegisterM2CConsumer() chan bool {
	c := make(chan bool)
	ssc.m2cCallbackConsumer = append(ssc.m2cCallbackConsumer, c)

	return c
}

// Registers a consumer of the module to the multifunction button
// Sends how long the button was held for
func (ssc *ShieldControl) RegisterMFBConsumer() chan uint16 {
	c := make(chan uint16)
	ssc.mfbCallbackConsumer = append(ssc.mfbCallbackConsumer, c)

	return c
}

// Function for checking if the external inputs were pressed and will signal consumers when ready
func (ssc *ShieldControl) btnCheck() {
	// if we have a signal from a downstream controller signal all consumers (nonblocking)
	if ssc.modInterruptPin.EdgeDetected() {
		for _, c := range ssc.m2cCallbackConsumer {
			go func(c chan bool) {
				c <- true
			}(c)
		}
	}
	// if the button state is changed, and we are looking for a press
	if ssc.mfbPin.EdgeDetected() && ssc.mfbEdge == rpio.RiseEdge {
		// Span a new concurent to wait for release
		go func() {
			// Set up the edge detection
			ssc.mfbEdge = rpio.FallEdge
			ssc.mfbPin.Detect(rpio.FallEdge)
			// Record time of press
			mfbPush := time.Now()
			// Wait for release
			for !ssc.mfbPin.EdgeDetected() {
			}
			// Record time of release and compute difference
			mfbRelease := time.Now()
			mfbPushTime := uint16(mfbRelease.Sub(mfbPush).Milliseconds())
			// Signal all consumers for how long the button was pressed
			for _, c := range ssc.mfbCallbackConsumer {
				// Non blocking channel update
				go func(c chan uint16) {
					c <- mfbPushTime
				}(c)
			}
			// Reset edge detection
			ssc.mfbEdge = rpio.RiseEdge
			ssc.mfbPin.Detect(rpio.RiseEdge)
		}()
	}
}
