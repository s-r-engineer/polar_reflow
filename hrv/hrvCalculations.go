package hrv

import (
	"math"
	influxclient "polar_reflow/influxClient"
	"polar_reflow/linker"
	"polar_reflow/syncronization"
	"polar_reflow/tools"
	"strings"
)

func RMSSD(rrIntervals []float64) float64 {
	if len(rrIntervals) < 2 {
		return 0
	}

	var sumSquares float64
	for i := 1; i < len(rrIntervals); i++ {
		diff := rrIntervals[i] - rrIntervals[i-1]
		sumSquares += diff * diff
	}
	rmssd := math.Sqrt(sumSquares / float64(len(rrIntervals)-1))
	return rmssd
}

func SDNN(pulseIntervals []float64) float64 {
	n := float64(len(pulseIntervals))
	if n == 0 {
		return 0.0
	}

	var sum float64 = 0
	for _, interval := range pulseIntervals {
		sum += interval
	}
	mean := sum / n

	var variance float64 = 0
	for _, interval := range pulseIntervals {
		variance += math.Pow(interval-mean, 2)
	}

	sdnn := math.Sqrt(variance / (n - 1))
	return sdnn
}

func hrvWorker(done func()) {
	defer done()
	for {
		optionInterface, _ := linker.Pop()
		if optionInterface == nil {
			break
		}
		option := optionInterface.([]string)
		method := option[0]
		timeTag := option[1]
		startTime := option[2]
		endTime := option[3]
		queryResult := influxclient.QueryPPI(startTime, endTime)

		result := []float64{}
		for queryResult.Next() {
			result = append(result, queryResult.Record().Value().(float64))
		}

		createHRVPoint(method, startTime, timeTag, result)
	}
}

func createHRVPoint(method string, startTimeS string, timeTag string, result []float64) {
	var data float64
	switch strings.ToLower(method) {
	case "sdnn":
		data = SDNN(result)
	case "rmssd":
		data = RMSSD(result)
	default:
		panic(method)
	}
	if data != 0.0 {
		startTime := tools.ParseTime(startTimeS)
		influxclient.WriteHRVPoint(timeTag, method, data, startTime)
	}
}

func SpinHRVWorkers(parallelism int) {
	add, done, wait := syncronization.CreateWGInstance()
	for range parallelism {
		add()
		go hrvWorker(done)
	}
	wait()
}
