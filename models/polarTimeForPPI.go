package models

import (
	"strings"
	"time"
)

const polarTimeFormatForPPI = `2006-01-02T15:04:05.000`

type polarTimeForPPI time.Time

func (pt *polarTimeForPPI) UnmarshalJSON(b []byte) error {
	s1 := strings.Split(strings.ReplaceAll(string(b), `"`, ""), "T")
	time1 := strings.Split(s1[1], ":")
	if len(time1) == 1 {
		time1 = append(time1, "00")
	}
	if len(time1) == 2 {
		time1 = append(time1, "00.000")
	}
	time1 = strings.Split(strings.Join(time1, ":"), ".")
	if len(time1) == 1 {
		time1 = append(time1, "000")
	}
	s1[1] = strings.Join(time1, ".")
	t, err := time.Parse(polarTimeFormatForPPI, strings.Join(s1, "T"))
	if err != nil {
		return err
	}

	*pt = polarTimeForPPI(t)
	return nil
}
