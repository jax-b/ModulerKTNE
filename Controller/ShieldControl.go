package Controller

import (
	"github.com/jax-b/go-i2c7Seg"
	"github.com/stianeikeland/go-rpio/v4"
)
struct ShieldControl {
	strikecount uint8
	seg *i2c7Seg.SevenSegI2C
	buzzerPin rpio.Pin
	strike1Pin rpio.Pin
	strike2Pin rpio.Pin
	modInterruptPin rpio.Pin
	randomStartPin rpio.Pin
}

func NewShieldControl(bus int, buzzerPinNum uint8, strike1PinNum uint8, strike2PinNum uint8, modInterruptPinNum uint8, randomStartPinNum uint8) (*ShieldControl, error) {
	i2c, err := i2c7Seg.NewSevenSegI2C(0x70, bus)
	if err != nil {
		return nil, err
	}
	err := rpio.Open()
	if err != nil {
		return nil, err
	}
	sc := &ShieldControl{
		seg: i2c,
		strikecount: 0,
		buzzerPin : rpio.Pin(buzzerPinNum),
		strike1Pin : rpio.Pin(strike1PinNum),
		strike2Pin : rpio.Pin(strike2PinNum),
		modInterruptPin : rpio.Pin(modInterruptPinNum),
		randomStartPin : rpio.Pin(randomStartPinNum),
	}
	// Configure pins
	sc.buzzerPin.Mode(rpio.Pwm)
	sc.strike1Pin.Output()
	sc.strike2Pin.Output()
	sc.randomStartPin.Input()
	sc.randomStartPin.PullUp()
	sc.randomStartPin.Detect(rpio.RiseEdge)
	sc.modInterruptPin.Input()
	sc.modInterruptPin.PullUp()
	sc.modInterruptPin.Detect(rpio.FallEdge)

	return sc, nil
}
func (self *ShieldControl) Close() {
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
	buzzerPin.Freq(64000)
	buzzerPin.DutyCycle(16, 32)
	time.Sleep(time.Millisecond * 500)
	buzzerPin.DutyCycle(0, 32)
}
// TimeClockSignal
func (self *ShieldControl) TimeSigBeep(multiplier float32) {
	buzzerPin.Freq(600)
	buzzerPin.DutyCycle(16, 32)
	time.Sleep(time.Millisecond * (350 * multiplier))
	buzzerPin.Freq(60)
	time.Sleep(time.Millisecond * (250 * multiplier)
	buzzerPin.DutyCycle(0, 32)
}
func (self *ShieldControl) WriteIdle() {
	self.seg.Clear()
	self.seg.WriteAsciiChar(0, 'K', false)
	self.seg.WriteAsciiChar(1, 'T', true)
	self.seg.WriteAsciiChar(3, 'N', false)
	self.seg.WriteAsciiChar(4, 'E', true)
	self.seg.WriteDisplay()
}
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
func (self *ShieldControl) ClearDisplay() {
	self.seg.Clear()
	self.seg.WriteDisplay()
}