package main

import (
	"fmt"
	"os"
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

	fmt.Println("Waiting for button press")
	for !randomStartPin.EdgeDetected() {
	}
	fmt.Println("button pressed")

	fmt.Println("1 Strike")
	strike1Pin.High()
	time.Sleep(time.Second * 1)
	fmt.Println("2 Strike")
	strike2Pin.High()
	time.Sleep(time.Second * 1)
	fmt.Println("3 Strike")
	strike1Pin.Low()
	time.Sleep(time.Second * 1)
	fmt.Println("No Strike")
	strike2Pin.Low()
	time.Sleep(time.Second * 1)

	if os.Getenv("EUID") == "0" {
		buzzerPin.Freq(64000)
		buzzerPin.DutyCycle(30, 32)
		time.Sleep(time.Millisecond * 500)
		buzzerPin.Freq(500)
		time.Sleep(time.Millisecond * 500)
		buzzerPin.DutyCycle(0, 32)
	} else {
		fmt.Println("not sudo: no buzzer")
	}

	fmt.Println("Setting up 7Seg I2C")
	sevenSeg, err := i2c7Seg.NewSevenSegI2C(0x70, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Writing KTNE")
	sevenSeg.WriteAsciiChar(0, 'K', false)
	sevenSeg.WriteAsciiChar(1, 'T', true)
	sevenSeg.WriteAsciiChar(3, 'N', false)
	sevenSeg.WriteAsciiChar(4, 'E', true)
	sevenSeg.WriteDisplay()
	time.Sleep(time.Second * 2)
	fmt.Println("Writing JAXB")
	sevenSeg.WriteAsciiChar(0, 'J', true)
	sevenSeg.WriteAsciiChar(1, 'A', false)
	sevenSeg.WriteAsciiChar(3, 'X', true)
	sevenSeg.WriteAsciiChar(4, 'B', false)
	sevenSeg.DrawColon(true)
	sevenSeg.WriteDisplay()
	time.Sleep(time.Second * 2)
	sevenSeg.Clear()
	sevenSeg.WriteDisplay()
	sevenSeg.Close()

	fmt.Println("Done")

}
