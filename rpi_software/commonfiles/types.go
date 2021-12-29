package commonfiles

import (
	"time"
)

type Status struct {
	StrTime             string `json:"timeleft"`
	Time                time.Duration
	NumStrike           uint8   `json:"strike"`
	Boom                bool    `json:"boom"`
	Win                 bool    `json:"win"`
	Gamerun             bool    `json:"gamerun"`
	Strikereductionrate float32 `json:"strikerate"`
}
