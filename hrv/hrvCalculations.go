package hrv

import (
	"math"
	"polar_reflow/database"
	influxclient "polar_reflow/database/influxClient"
	"polar_reflow/linker"
	"polar_reflow/models"
	"polar_reflow/syncronization"
	"polar_reflow/tools"
	"reflect"
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
		result := database.Get(tools.ParseTime(startTime), tools.ParseTime(endTime), nil)
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
	//timePoint := time.Unix(int64(t/1000), 0)
	//result := database.Get(timePoint.Add(time.Minute*-5), timePoint)
	//return RMSSD(result)
	return 0
}

func Get5MinRMSSDFromtimeToTime(t1, t2 time.Time) (resultPoints []models.DBPPI) {
	minutesChannel := make(chan []models.DBPPI, 10000)
	resultChannel := [16]chan models.DBPPI{}
	add, done, wait := syncronization.CreateWGInstance()
	for i := 0; i < 16; i++ {
		c := i
		resultChannel[c] = make(chan models.DBPPI, 6000)
		go func() {
			add()
			defer done()
			for {
				readyMinute, ok := <-minutesChannel
				if !ok {
					close(resultChannel[c])
					return
				}
				if CheckAmountOfPoints(readyMinute, time.Duration(5*time.Minute)) {
					resultChannel[c] <- models.DBPPI{Value: RMSSD(readyMinute), TimePoint: readyMinute[len(readyMinute)-1].TimePoint}
				}
			}
		}()
	}
	go func() {
		add()
		defer done()
		adder := func(result models.DBPPI) {
			resultPoints = append(resultPoints, result)
		}
		var cases []reflect.SelectCase
		for _, ch := range resultChannel {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}

		for len(cases) > 0 {
			chosen, value, ok := reflect.Select(cases)
			//tools.Dumper(chosen)
			if !ok {
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			//tools.Dumper(value.Interface().(models.DBPPI))
			adder(value.Interface().(models.DBPPI))
		}
	}()
	//for {
	//	select {
	//	case result := <-resultChannel[0]:
	//	case result := <-resultChannel[2]:
	//	case result := <-resultChannel[3]:
	//	case result := <-resultChannel[4]:
	//	case result := <-resultChannel[5]:
	//	case result := <-resultChannel[6]:
	//	case result := <-resultChannel[7]:
	//	case result := <-resultChannel[8]:
	//	case result := <-resultChannel[9]:
	//	case result := <-resultChannel[10]:
	//	case result := <-resultChannel[11]:
	//	case result := <-resultChannel[12]:
	//	case result := <-resultChannel[13]:
	//	case result := <-resultChannel[14]:
	//	case result := <-resultChannel[15]:
	//		adder(result)
	//	}
	//}

	channel := make(chan models.DBPPI, 700000)
	go database.Get(t1, t2, channel)
	minutes := make(map[int][]models.DBPPI)

	// TODO fix this. Point must be counted
	val, ok := <-channel
	if !ok {
		return
	}
	//aqq, release := syncronization.CreateMutexInstance()
	startPoint := val.TimePoint
	for {
		val, ok = <-channel
		if ok {
			currentPoint := int(val.TimePoint.Sub(startPoint).Minutes())
			//aqq()
			for m := currentPoint + 1; m <= currentPoint+5; m++ {
				if _, ok := minutes[m]; !ok {
					minutes[m] = []models.DBPPI{val}
				} else {
					minutes[m] = append(minutes[m], val)
				}
			}

			readyMinute, ok := minutes[currentPoint]
			if ok {
				//logger.Debug(fmt.Sprintf("spinning minute counter for minute %d", currentPoint))
				minutesChannel <- readyMinute
				delete(minutes, currentPoint)
			}

			//release()
		} else {
			close(minutesChannel)
			break
		}
	}
	wait()
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
