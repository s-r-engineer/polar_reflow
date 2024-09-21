package hrv

import (
	"math"
	influxclient "polar_reflow/influxClient"
	"polar_reflow/linker"
	"polar_reflow/models"
	"polar_reflow/syncronization"
	"polar_reflow/tools"
	"sort"
	"strings"
	"time"
)

func RMSSD(rrIntervals []models.PPIFromInflux) float64 {
	if len(rrIntervals) < 2 {
		return 0
	}

	var sumSquares float64
	for i := 1; i < len(rrIntervals); i++ {
		diff := rrIntervals[i].Value - rrIntervals[i-1].Value
		sumSquares += diff * diff
	}
	rmssd := math.Sqrt(sumSquares / float64(len(rrIntervals)-1))
	return rmssd
}

func SDNN(pulseIntervals []models.PPIFromInflux) float64 {
	n := float64(len(pulseIntervals))
	if n == 0 {
		return 0.0
	}

	var sum float64 = 0
	for _, point := range pulseIntervals {
		sum += point.Value
	}
	mean := sum / n

	var variance float64 = 0
	for _, point := range pulseIntervals {
		variance += math.Pow(point.Value-mean, 2)
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
		result := getValuesFromInflux(startTime, endTime)
		createHRVPoint(method, startTime, timeTag, result)
	}
}

func getValuesFromInflux(startTime, endTime string) (result []models.PPIFromInflux) {
	queryResult := influxclient.QueryPPI(startTime, endTime)
	for queryResult.Next() {
		result = append(result, models.PPIFromInflux{Value: queryResult.Record().Value().(float64), TimePoint: queryResult.Record().Time()})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TimePoint.Before(result[j].TimePoint)
	})
	return
}

func createHRVPoint(method string, startTimeS string, timeTag string, result []models.PPIFromInflux) {
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

func Get5MinRMSSDFromPoint(t int) float64 {
	timePoint := time.Unix(int64(t/1000), 0)
	result := getValuesFromInflux(tools.FormatTime(timePoint.Add(time.Minute*-5)), tools.FormatTime(timePoint))
	return RMSSD(result)
}

func Get5MinRMSSDFromtimeToTime(t1, t2 time.Time) []models.PPIFromInflux {
	result := getValuesFromInflux(tools.FormatTime(t1), tools.FormatTime(t2))
	minutes := make(map[int][]models.PPIFromInflux)
	resultPoints := []models.PPIFromInflux{}
	if len(result) == 0 {
		return resultPoints
	}
	startPoint := result[0].TimePoint

	for _, pointInResult := range result {
		currentPoint := int(pointInResult.TimePoint.Sub(startPoint).Minutes())
		for m := currentPoint; m > currentPoint-5 && m > 0; m-- {
			if _, ok := minutes[m]; !ok {
				minutes[m] = []models.PPIFromInflux{}
			}
			minutes[m] = append(minutes[m], pointInResult)
		}
	}
	for _, v := range minutes {
		r := RMSSD(v)
		resultPoints = append(resultPoints, models.PPIFromInflux{Value: r, TimePoint: v[len(v)-1].TimePoint})
	}
	tools.Dumper(resultPoints[0])

	return resultPoints
}
