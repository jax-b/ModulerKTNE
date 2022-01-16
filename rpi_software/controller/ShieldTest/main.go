package main

import (
	"log"
	"time"
	"os"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/jax-b/go-i2c7Seg"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	log.Println("Hello world")
	const strike1PinNum uint8 = 22
	const strike2PinNum uint8 = 23
	// const modInterruptPinNum uint8 = 17
	const randomStartPinNum uint8 = 27

	err := rpio.Open()
	if err != nil {
		log.Println(err)
	}

	defer rpio.Close()

	strike1Pin := rpio.Pin(strike1PinNum)
	strike2Pin := rpio.Pin(strike2PinNum)
	// modInterruptPin := rpio.Pin(modInterruptPinNum)
	randomStartPin := rpio.Pin(randomStartPinNum)

	strike1Pin.Output()
	strike2Pin.Output()

	randomStartPin.Input()
	randomStartPin.PullUp()
	randomStartPin.Detect(rpio.FallEdge)

	log.Println("Waiting for button press")
	for !randomStartPin.EdgeDetected() {
	}
	log.Println("button pressed")
	randomStartPin.Detect(rpio.NoEdge)

	log.Println("1 Strike")
	strike1Pin.High()
	time.Sleep(time.Second * 1)
	log.Println("2 Strike")
	strike2Pin.High()
	time.Sleep(time.Second * 1)
	log.Println("3 Strike")
	strike1Pin.Low()
	time.Sleep(time.Second * 1)
	log.Println("No Strike")
	strike2Pin.Low()
	time.Sleep(time.Second * 1)

	log.Println("Setting up 7Seg I2C")
	sevenSeg, err := i2c7Seg.NewSevenSegI2C(0x70, 1)
	if err != nil {
		log.Println(err)
	}
	log.Println("Writing KTNE")
	sevenSeg.WriteAsciiChar(0, 'K', false)
	sevenSeg.WriteAsciiChar(1, 'T', true)
	sevenSeg.WriteAsciiChar(3, 'N', false)
	sevenSeg.WriteAsciiChar(4, 'E', true)
	sevenSeg.WriteDisplay()
	time.Sleep(time.Second * 2)
	log.Println("Writing JAXB")
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

	log.Println("Sound Test")
	soundfile, err := os.Open("./audiotst.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(soundfile)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	sr := format.SampleRate * 2
	log.Println(sr)
	speaker.Init(sr, sr.N(time.Second/10))

	resampled := beep.Resample(4, format.SampleRate, sr, streamer)
	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	<-done

	log.Println("Done")

}
