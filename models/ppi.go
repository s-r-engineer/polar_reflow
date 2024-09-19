package models

import "time"

type PPI []struct {
	Date                 string                 `json:"date"`
	DevicePpiSamplesList []DevicePpiSamplesList `json:"devicePpiSamplesList"`
}

type DevicePpiSamplesList struct {
	DeviceID   string       `json:"deviceId"`
	PpiSamples []PpiSamples `json:"ppiSamples"`
}

type PpiSamples struct {
	SampleDateTime polarTime `json:"sampleDateTime"`
	TimeWithPulse  time.Time
	PulseLength    int `json:"pulseLength"`
}
