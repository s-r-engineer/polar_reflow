package main

import (
	"flag"
	"fmt"
	"polar_reflow/hrv"
	importData "polar_reflow/import"
	influxclient "polar_reflow/influxClient"
	"polar_reflow/linker"
	mygin "polar_reflow/myGin"
	"polar_reflow/tools"
	"time"
)

var (
	pathToFindFilesIn         = flag.String("path", "", "")
	reinit                    = flag.Bool("reinit", false, "")
	influxAddress             = flag.String("influx", "http://localhost:8086", "")
	excludeRmssd              = flag.Bool("normssd", false, "")
	excludeSddn               = flag.Bool("nosddn", false, "")
	token                     = flag.String("token", "", "")
	org                       = flag.String("org", "my", "")
	bucket                    = flag.String("bucket", "user", "")
	startTimeString           = flag.String("start", "2020-01-01T00:00:00Z", "")
	finaltime                 = flag.String("end", time.Now().Format("2006-01-02T15:04:05Z"), "")
	parallelismForCalculating = flag.Int("paralel", 16, "")

	periods = map[string][]int{"sdnn": {2 * 60, 12 * 60, 24 * 60}, "rmssd": {5}}
)

func main() {
	flag.Parse()
	go mygin.Run()
	influxclient.InitInflux(*influxAddress, *token, *org, *bucket)
	if *reinit {
		influxclient.ReinitBucket()
		return
	}
	if *pathToFindFilesIn != "" {
		importData.ImportFiles(*pathToFindFilesIn)
	}

	fmt.Println("Starting calculations")

	finalTimeO, err := time.Parse("2006-01-02T15:04:05Z07:00", *finaltime)
	tools.ErrPanic(err)

	startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", *startTimeString)
	tools.ErrPanic(err)

	linker.CreateLinker(*excludeSddn, *excludeRmssd, startTime, finalTimeO, periods)
	hrv.SpinHRVWorkers(*parallelismForCalculating)
	influxclient.Flush()

	fmt.Println("Done calculations")
}
