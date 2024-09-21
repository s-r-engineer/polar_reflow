package tools

import (
	"time"

	"github.com/davecgh/go-spew/spew"
)

const customTimeFormat = "2006-01-02T15:04:05"

func ErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func FormatTime(t time.Time) string {
	return t.Format(customTimeFormat)
}

func ParseTime(t string) time.Time {
	timePoint, err := time.Parse(customTimeFormat, t)
	ErrPanic(err)
	return timePoint
}

func Dumper(v any) {
	spew.Dump(v)
}
