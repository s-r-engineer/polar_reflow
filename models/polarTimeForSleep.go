package models

import (
	"strings"
	"time"
)

const polarTimeFormatForSleep = `2006-01-02`

type PolarTimeForSleep time.Time

func (pt *PolarTimeForSleep) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(polarTimeFormatForSleep, strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}

	*pt = PolarTimeForSleep(t)
	return nil
}
