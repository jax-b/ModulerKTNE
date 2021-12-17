package netdisplay

import "fmt"

func timetostring(cmiltime uint32) string {
	var tstring string
	timemilisf := float32(cmiltime)
	MinutesRemaining := int(cmiltime*0.001) / 60
	SecondsRemaining := int(cmiltime*0.001) % 60
	if MinutesRemaining > 0 {
		if MinutesRemaining > 99 {
			MinutesRemaining = 99
		}
		tstring = fmt.Sprintf("%2d:%02d", MinutesRemaining, SecondsRemaining)
	} else {
		hundtensec := int(timemilisf*0.1) % 100
		tstring = fmt.Sprintf("%2d.%02d", SecondsRemaining, hundtensec)
	}
	return tstring
}
