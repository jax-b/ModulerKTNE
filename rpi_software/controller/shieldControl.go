package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/jax-b/go-i2c7Seg"
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/zap"
)

type ShieldControl struct {
	strikecount uint8
	seg         *i2c7Seg.SevenSegI2C

	strike1Pin          rpio.Pin
	strike2Pin          rpio.Pin
	m2cPin              rpio.Pin
	mfbPin              rpio.Pin
	mfbCallbackConsumer []chan time.Duration
	m2cCallbackConsumer []chan bool
	stopM2CCheck        chan bool
	stopMFBCheck        chan bool
	log                 *zap.SugaredLogger
	samplerate          beep.SampleRate
}

func NewShieldControl(logger *zap.SugaredLogger, cfg *Config) *ShieldControl {
	logger = logger.Named("ShieldControl")
	logger.Info("Starting Shield Control")

	i2c, err := i2c7Seg.NewSevenSegI2C(cfg.Shield.SevenSegAddress, int(cfg.Shield.I2cBusNumber))
	i2c.LogLevel("INFO")
	if err != nil {
		logger.Error("Failed to create i2c7Seg", err)
	}
	err = rpio.Open()
	if err != nil {
		logger.Error("Failed to open rpio", err)
	}
	sc := &ShieldControl{
		seg:         i2c,
		strikecount: 0,
		strike1Pin:  rpio.Pin(cfg.Shield.Strike1PinNum),
		strike2Pin:  rpio.Pin(cfg.Shield.Strike2PinNum),
		m2cPin:      rpio.Pin(cfg.Shield.ModInterruptPinNum),
		mfbPin:      rpio.Pin(cfg.Shield.MfbStartPinNum),
		log:         logger,
	}

	// Configure pins
	logger.Info("Configuring Output Pins")
	sc.strike1Pin.Output()
	sc.strike1Pin.Low()
	sc.strike2Pin.Output()
	sc.strike2Pin.Low()

	sc.mfbPin.Input()
	sc.mfbPin.PullUp()
	sc.mfbPin.Detect(rpio.FallEdge)

	sc.m2cPin.Input()
	sc.m2cPin.PullUp()
	sc.m2cPin.Detect(rpio.FallEdge)

	// Configure Speaker
	logger.Info("Configuring Output Speaker")
	sc.samplerate = beep.SampleRate(48000)
	speaker.Init(sc.samplerate, sc.samplerate.N(time.Second/10))

	sc.stopMFBCheck = make(chan bool, 1) // Buffered channel so we dont hang on close if run is never called
	sc.stopM2CCheck = make(chan bool, 1) // Buffered channel so we dont hang on close if run is never called
	return sc
}

func (ssc *ShieldControl) Run() {
	// Start Input checker
	go ssc.m2cCheck()
	go ssc.mfbCheck()
}

// Closes out all functions that are running safely
func (ssc *ShieldControl) Close() {
	ssc.stopMFBCheck <- true
	ssc.stopM2CCheck <- true
	time.After(50 * time.Millisecond)

	ssc.mfbPin.Detect(rpio.NoEdge)
	ssc.m2cPin.Detect(rpio.NoEdge)
	ssc.ClearDisplay()
	ssc.seg.Close()
	ssc.ResetStrike()
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

// Adds a strike to the display system and plays a sound
func (ssc *ShieldControl) SetStrike(numstrike uint8) {
	ssc.strikecount = numstrike
	// span strike sound in a seprate concurent
	if ssc.strikecount == 0 {
		ssc.strike1Pin.Low()
		ssc.strike2Pin.Low()
	} else if ssc.strikecount == 1 {
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
	// Open The Audio file
	soundfile, err := os.Open("./audiofiles/strike.wav")
	if err != nil {
		ssc.log.Error(err)
		return
	}

	// Decode its format
	streamer, format, err := wav.Decode(soundfile)
	if err != nil {
		ssc.log.Error(err)
		return
	}
	defer streamer.Close()

	// Resample it to our standered for output
	resampled := beep.Resample(4, format.SampleRate, ssc.samplerate, streamer)

	// Play the file and wait for it to finish
	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}

// plays the sound of an explosion
func (ssc *ShieldControl) ExploadSound() {
	// Open The Audio file
	explof, err := os.Open("./audiofiles/explosion_concrete_medium.wav")
	if err != nil {
		ssc.log.Error(err)
		return
	}

	// Decode its format
	explostreamer, format, err := wav.Decode(explof)
	if err != nil {
		ssc.log.Error(err)
		return
	}
	defer explostreamer.Close()

	// Resample it to our standered for output
	esxploresampled := beep.Resample(4, format.SampleRate, ssc.samplerate, explostreamer)

	// Play the file and wait for it to finish
	done := make(chan bool)
	speaker.Play(beep.Seq(esxploresampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}

// plays the sound of an game win
func (ssc *ShieldControl) GameWinSound() {
	// Open The Audio file
	winsoundf, err := os.Open("./audiofiles/mktne-winmix.wav")
	if err != nil {
		ssc.log.Error(err)
		return
	}

	// Decode its format
	winsoundstream, format, err := wav.Decode(winsoundf)
	if err != nil {
		ssc.log.Error(err)
		return
	}
	defer winsoundstream.Close()

	// Resample it to our standered for output
	winsoundsampled := beep.Resample(4, format.SampleRate, ssc.samplerate, winsoundstream)

	// Play the file and wait for it to finish
	done := make(chan bool)
	speaker.Play(beep.Seq(winsoundsampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}

// plays the sound of an needy wanting attention
func (ssc *ShieldControl) NeedyWantSound() {
	// Open The Audio file
	needyf, err := os.Open("./audiofiles/needy_activated.wav")
	if err != nil {
		ssc.log.Error(err)
		return
	}

	// Decode its format
	needystream, format, err := wav.Decode(needyf)
	if err != nil {
		ssc.log.Error(err)
		return
	}
	defer needystream.Close()

	// Resample it to our standered for output
	needysampled := beep.Resample(4, format.SampleRate, ssc.samplerate, needystream)

	// Play the file and wait for it to finish
	done := make(chan bool)
	speaker.Play(beep.Seq(needysampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}

// plays the sound of an module that has been solved
func (ssc *ShieldControl) ModSolvedSound() {
	// Open The Audio file
	modcorf, err := os.Open("./audiofiles/CorrectDigitalChime.wav")
	if err != nil {
		ssc.log.Error(err)
		return
	}

	// Decode its format
	modcorstr, format, err := wav.Decode(modcorf)
	if err != nil {
		ssc.log.Error(err)
		return
	}
	defer modcorstr.Close()

	// Resample it to our standered for output
	modcorresamp := beep.Resample(4, format.SampleRate, ssc.samplerate, modcorstr)

	// Play the file and wait for it to finish
	done := make(chan bool)
	speaker.Play(beep.Seq(modcorresamp, beep.Callback(func() {
		done <- true
	})))

	<-done
}

// TimeClockSignal Plays the per second beep
// takes the channel needed to stop the beep and a tick chan to trigger each tone on the game clock
func (ssc *ShieldControl) TimerBeep(stopchannel chan bool, timertick chan bool) {
	go func(stopch chan bool, tmrTick chan bool) {
		// Open The Audio file
		timeBeepf, err := os.Open("./audiofiles/doublebeep.wav")
		if err != nil {
			ssc.log.Error("Entering Dead State", err)
			// if we have a error we dont want the program to lock because of a waiting chan
			// so we just constently read each chan and and wait to exit
			go func(stopch chan bool, tmrTick chan bool) {
				select {
				case <-stopch:
					return
				case <-tmrTick:
				}
			}(stopchannel, timertick)
			return
		}

		timeBeepStreamer, format, err := wav.Decode(timeBeepf)
		if err != nil {
			ssc.log.Error("Entering Dead State", err)
			go func(stopch chan bool, tmrTick chan bool) {
				// if we have a error we dont want the program to lock because of a waiting chan
				// so we just constently read each chan and and wait to exit
				select {
				case <-stopch:
					return
				case <-tmrTick:
				}
			}(stopchannel, timertick)
			return
		}
		// Resample it to our standered for output
		timeBeepResampled := beep.Resample(4, format.SampleRate, ssc.samplerate, timeBeepStreamer)

		//Create a buffer in memory to hold the audio file
		buffer := beep.NewBuffer(beep.Format{
			Precision:   format.Precision,
			SampleRate:  ssc.samplerate,
			NumChannels: format.NumChannels,
		})
		// load the audio file into the buffer
		buffer.Append(timeBeepResampled)
		// close the file and streamer
		timeBeepStreamer.Close()

		for { // Wait to either close the playback system or for a tick from the timer to play the sound file
			select {
			case <-stopch:

				return
			case <-tmrTick:
				tone := buffer.Streamer(0, buffer.Len()) // Load a streamer from the zero spot in the file
				speaker.Play(tone)                       // play the loaded streamer
			}
		}

	}(stopchannel, timertick)
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
func (ssc *ShieldControl) WriteTime(intime time.Duration) {
	ssc.seg.Clear()
	var tstring string
	MinutesRemaining := int(intime.Minutes())
	SecondsRemaining := int(intime.Seconds()) % 60
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
		hundtensec := intime.Milliseconds() / 10 % 100
		// fmt.Println(hundtensec)
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
func (ssc *ShieldControl) RegisterMFBConsumer() chan time.Duration {
	c := make(chan time.Duration)
	ssc.mfbCallbackConsumer = append(ssc.mfbCallbackConsumer, c)

	return c
}

// Function for checking if the external inputs were pressed and will signal consumers when ready
func (ssc *ShieldControl) mfbCheck() {
	mfbEdge := rpio.FallEdge
	for {
		select {
		case <-ssc.stopMFBCheck:
			return
		default:
		}
		// if the button state is changed, and we are looking for a press
		if mfbEdge == rpio.FallEdge {
			if ssc.mfbPin.EdgeDetected() {
				ssc.log.Info("Detected Button Press, Waiting for Release")
				// Set up the edge detection
				mfbEdge = rpio.RiseEdge
				ssc.mfbPin.Detect(mfbEdge)
				// Record time of press
				mfbPush := time.Now()
				// Span a new concurent to wait for release
				go func(pushtime time.Time) {
					// Wait for release
					timeMinExceded := false
					for !timeMinExceded {
						for !ssc.mfbPin.EdgeDetected() {
						}
						if time.Now().Sub(pushtime) > time.Millisecond*4 {
							timeMinExceded = true
						}
					}

					// Record time of release and compute difference
					mfbHeldTime := time.Now().Sub(pushtime)
					ssc.log.Info("Button Release Detected, Notifiying Consumers of Duration: ", mfbHeldTime)

					// Signal all consumers for how long the button was pressed
					for _, c := range ssc.mfbCallbackConsumer {
						// Non blocking channel update
						go func(c chan time.Duration) {
							c <- mfbHeldTime
						}(c)
					}

					// Reset edge detection
					mfbEdge = rpio.FallEdge
					ssc.mfbPin.Detect(mfbEdge)
				}(mfbPush)
			}
		}
	}
}

func (ssc *ShieldControl) m2cCheck() {
	for {
		select {
		case <-ssc.stopM2CCheck:
			return
		default:
		}
		// if we have a signal from a downstream controller signal all consumers (nonblocking)
		if ssc.m2cPin.EdgeDetected() {
			ssc.log.Debug("M2C Interrupt Detected, Notifiying Consumers")
			for _, c := range ssc.m2cCallbackConsumer {
				go func(c chan bool) {
					c <- true
				}(c)
			}
			ssc.m2cPin.Detect(rpio.FallEdge)
			time.Sleep(5 * time.Millisecond)
		}
	}
}
