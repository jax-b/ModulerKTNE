package main

import (
	"fmt"
	"time"

	"github.com/jax-b/go-i2c7Seg"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	fmt.Println("Hello world")
	const buzzerPinNum uint8 = 18
	const strike1PinNum uint8 = 22
	const strike2PinNum uint8 = 23
	const modInterruptPinNum uint8 = 17
	const randomStartPinNum uint8 = 27

	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
	}

	defer rpio.Close()

	buzzerPin := rpio.Pin(buzzerPinNum)
	strike1Pin := rpio.Pin(strike1PinNum)
	strike2Pin := rpio.Pin(strike2PinNum)
	// modInterruptPin := rpio.Pin(modInterruptPinNum)
	randomStartPin := rpio.Pin(randomStartPinNum)

	buzzerPin.Mode(rpio.Pwm)
	strike1Pin.Output()
	strike2Pin.Output()

	randomStartPin.Input()
	randomStartPin.PullUp()
	randomStartPin.Detect(rpio.FallEdge)

	for !randomStartPin.EdgeDetected() {
	}
	fmt.Println("button pressed")

	strike1Pin.High()
	time.Sleep(time.Second * 1)
	strike2Pin.High()
	time.Sleep(time.Second * 1)
	strike1Pin.Low()
	time.Sleep(time.Second * 1)
	strike2Pin.Low()
	time.Sleep(time.Second * 1)

	buzzerPin.Freq(64000)
	buzzerPin.DutyCycle(30, 32)
	time.Sleep(time.Millisecond * 500)
	buzzerPin.Freq(500)
	time.Sleep(time.Millisecond * 500)
	buzzerPin.DutyCycle(0, 32)

	sevenSeg, err := i2c7Seg.NewSevenSegI2C(0x70, 1)
	if err != nil {
		fmt.Println(err)
	}
	sevenSeg.WriteAsciiChar(0, 'D', false)
	sevenSeg.WriteAsciiChar(1, 'E', true)
	sevenSeg.WriteAsciiChar(3, 'A', false)
	sevenSeg.WriteAsciiChar(4, 'D', true)
	sevenSeg.WriteDisplay()
	time.Sleep(time.Second * 2)
	sevenSeg.WriteAsciiChar(0, 'B', true)
	sevenSeg.WriteAsciiChar(1, 'E', false)
	sevenSeg.WriteAsciiChar(3, 'E', true)
	sevenSeg.WriteAsciiChar(4, 'F', false)
	sevenSeg.DrawColon(true)
	sevenSeg.WriteDisplay()
	time.Sleep(time.Second * 2)
	sevenSeg.Clear()
	sevenSeg.WriteDisplay()
}
