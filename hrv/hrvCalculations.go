package hrv

import (
	"math"
	"polar_reflow/database"
	influxclient "polar_reflow/database/influxClient"
	"polar_reflow/linker"
	"polar_reflow/models"
	"polar_reflow/syncronization"
	"polar_reflow/tools"
	"sort"
	"strings"
	"time"
)

func RMSSD(rrIntervals []models.DBPPI) float64 {
	if len(rrIntervals) < 2 {
		return 0
	}

	var sumSquares float64
	for i := 1; i < len(rrIntervals); i++ {
		diff := rrIntervals[i].Value - rrIntervals[i-1].Value
		sumSquares += diff * diff
	}

	rmssd := math.Sqrt(float64(sumSquares) / float64(len(rrIntervals)-1))
	return rmssd
}

func SDNN(pulseIntervals []models.DBPPI) float64 {
	n := float64(len(pulseIntervals))
	if n == 0 {
		return 0.0
	}

	var sum float64
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

func hrvWorker(done func(), pop func() (any, func())) {
	defer done()
	for {
		optionInterface, _ := pop()
		if optionInterface == nil {
			break
		}
		option := optionInterface.([]string)
		method := option[0]
		timeTag := option[1]
		startTime := option[2]
		endTime := option[3]
		result := database.GetWithTimeAsString(startTime, endTime)
		createHRVPoint(method, startTime, timeTag, result)
	}
}

func createHRVPoint(method string, startTimeS string, timeTag string, result []models.DBPPI) {
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

func SpinHRVWorkers(parallelism int, linker *linker.Linker) {
	add, done, wait := syncronization.CreateWGInstance()
	for range parallelism {
		add()
		go hrvWorker(done, linker.Pop)
	}
	wait()
}

func Get5MinRMSSDFromPoint(t int) float64 {
	timePoint := time.Unix(int64(t/1000), 0)
	result := database.Get(timePoint.Add(time.Minute*-5), timePoint)
	return RMSSD(result)
}

func Get5MinRMSSDFromtimeToTime(t1, t2 time.Time) []models.DBPPI {
	result := database.Get(t1, t2)
	resultPoints := []models.DBPPI{}
	if len(result) == 0 {
		return resultPoints
	}
	minutes := make(map[int][]models.DBPPI)
	startPoint := result[0].TimePoint
	for _, pointInResult := range result {
		currentPoint := int(pointInResult.TimePoint.Sub(startPoint).Minutes())
		for m := currentPoint + 1; m <= currentPoint+5; m++ {
			if _, ok := minutes[m]; !ok {
				minutes[m] = []models.DBPPI{}
			}
			minutes[m] = append(minutes[m], pointInResult)
		}
	}

	for _, v := range minutes {
		if CheckAmountOfPoints(v, time.Duration(5*time.Minute)) {
			r := RMSSD(v)
			resultPoints = append(resultPoints, models.DBPPI{Value: r, TimePoint: v[len(v)-1].TimePoint})
		}
	}
	sort.Slice(resultPoints, func(i, j int) bool {
		return resultPoints[i].TimePoint.Before(resultPoints[j].TimePoint)
	})
	return resultPoints
}

func CheckAmountOfPoints(rrIntervals []models.DBPPI, duration time.Duration) bool {
	for i := range rrIntervals {
		duration -= time.Duration(time.Millisecond * time.Duration(rrIntervals[i].Value))
	}
	return duration < time.Second*20
}
