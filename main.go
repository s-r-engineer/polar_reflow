package main

import (
	"flag"
	"os"
	"polar_reflow/hrv"

	// importData "polar_reflow/import"
	"polar_reflow/linker"
	"polar_reflow/logger"
	mygin "polar_reflow/myGin"
	"time"
)

var (
	pathToFindFilesIn         = flag.String("path", os.Getenv("POLAR_REFLOW_PATH"), "")
	serve                     = flag.Bool("serve", "true" == os.Getenv("POLAR_REFLOW_SERVE"), "")
	database                  = flag.String("database", os.Getenv("POLAR_REFLOW_DATABASE"), "Database backend. Could be mongo or influx")
	reinit                    = flag.Bool("reinit", "true" == os.Getenv("POLAR_REFLOW_REINIT"), "")
	influxAddress             = flag.String("influx_address", "http://influx:8086", "")
	influxAddress             = flag.String("influx", "http://influx:8086", "")
	excludeRmssd              = flag.Bool("normssd", "true" == os.Getenv("POLAR_REFLOW_NORMSSD"), "")
	excludeSddn               = flag.Bool("nosddn", "true" == os.Getenv("POLAR_REFLOW_NOSDDN"), "")
	token                     = flag.String("token", os.Getenv("POLAR_REFLOW_TOKEN"), "")
	org                       = flag.String("org", os.Getenv("POLAR_REFLOW_ORG"), "")
	bucket                    = flag.String("bucket", os.Getenv("POLAR_REFLOW_BUCKET"), "")
	startTimeString           = flag.String("start", "2020-01-01T00:00:00Z", "")
	finaltime                 = flag.String("end", time.Now().Format("2006-01-02T15:04:05Z"), "")
	parallelismForCalculating = flag.Int("paralel", 16, "")

	periods = map[string][]int{"sdnn": {2 * 60, 12 * 60, 24 * 60}, "rmssd": {5}}
)

func init() {
	logger.InitLogger(os.Getenv("POLAR_REFLOW_DEPLOYMENT"))
}

func main() {
	logger.Info("Starting")
	flag.Parse()
	// if *reinit {
	// 	influxclient.ReinitBucket()
	// 	return
	// }
	// if *pathToFindFilesIn != "" {
	// 	importData.ImportFiles(*pathToFindFilesIn)
	// 	return
	// }
	if *serve {
		mygin.Run()
		return
	}

	logger.Info("Starting calculations")

	finalTimeO, err := time.Parse("2006-01-02T15:04:05Z07:00", *finaltime)
	logger.Error(err.Error())

	startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", *startTimeString)
	logger.Error(err.Error())

	hrv.SpinHRVWorkers(*parallelismForCalculating, CreateLinkerForPeriods(*excludeSddn, *excludeRmssd, startTime, finalTimeO, periods))
	// influxclient.Flush()

	logger.Info("Done calculations")
}

func CreateLinkerForPeriods(excludeSddn, excludeRmssd bool, startTime, endTime time.Time, periods map[string][]int) *linker.Linker {
	linker := linker.Linker{}
	for method, timePeriods := range periods {
		if (excludeSddn && method == "sddn") || (excludeRmssd && method == "rmssd") {
			continue
		}
		for _, timePeriod := range timePeriods {
			clearHours := timePeriod / 60
			minutesLeft := timePeriod % 60
			clearDays := clearHours / 24
			hoursLeft := clearHours % 24
			offset := time.Duration(timePeriod) * time.Minute
			timeTagLine := fmt.Sprintf("%d%s%d%s%d%s", clearDays, "d", hoursLeft, "h", minutesLeft, "m")
			for timeCounter := startTime; timeCounter.Before(endTime); timeCounter = timeCounter.Add(offset) {
				linker.Push([]string{
					method, timeTagLine, tools.FormatTime(timeCounter), tools.FormatTime(timeCounter.Add(offset)),
				})
			}
		}
	}

	return &linker
}
