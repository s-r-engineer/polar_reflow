package main

import (
	"flag"
	"fmt"
	"os"
	"polar_reflow/hrv"
	importData "polar_reflow/import"
	influxclient "polar_reflow/influxClient"
	"polar_reflow/linker"
	mygin "polar_reflow/myGin"
	"polar_reflow/tools"
	"time"
)

var (
	pathToFindFilesIn         = flag.String("path", os.Getenv("POLAR_REFLOW_PATH"), "")
	serve                     = flag.Bool("serve", "true" == os.Getenv("POLAR_REFLOW_SERVE"), "")
	reinit                    = flag.Bool("reinit", "true" == os.Getenv("POLAR_REFLOW_REINIT"), "")
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

func main() {
	flag.Parse()
	influxclient.InitInflux(*influxAddress, *token, *org, *bucket)
	if *reinit {
		influxclient.ReinitBucket()
		return
	}
	if *pathToFindFilesIn != "" {
		importData.ImportFiles(*pathToFindFilesIn)
		return
	}
	if *serve {
		mygin.Run()
		return
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
