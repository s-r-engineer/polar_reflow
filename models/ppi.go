package models

import (
	"github.com/google/uuid"
	"time"
)

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

type DBPPI struct {
	Value     float64
	TimePoint time.Time
}

type BTHR struct {
	Session   uuid.UUID
	Value     uint16
	TimePoint time.Time
}

type MongoDBPPI struct {
	Value     float64   `bson:"value"`
	TimePoint time.Time `bson:"timePoint"`
	ID        string    `bson:"_id"`
}
