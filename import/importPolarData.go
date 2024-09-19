package importData

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	influxclient "polar_reflow/influxClient"
	"polar_reflow/models"
	"polar_reflow/syncronization"
	"polar_reflow/tools"
	"regexp"
	"time"
)

func ImportFiles(pathToLookIn string) {
	fmt.Println("Starting reading files")
	absPath, err := filepath.Abs(pathToLookIn)
	tools.ErrPanic(err)
	aqquire, release := syncronization.CreateSemaphoreInstance(4)
	add, done, wait := syncronization.CreateWGInstance()
	tools.ErrPanic(filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		return importFile(path, aqquire, release, add, done)
	},
	))
	wait()
	fmt.Println("flushing data")

	influxclient.Flush()
}

func importFile(path string, aqquire func() error, release, add, done func()) error {
	m, err := regexp.MatchString(`^.*ppi_.*\.json$`, filepath.Base(path))
	tools.ErrPanic(err)
	if !m {
		return nil
	}
	add()
	go func(path string) {
		defer done()
		tools.ErrPanic(aqquire())
		tools.ErrPanic(err)

		defer release()

		fmt.Printf("file %s parsing\n", path)

		reader, err := os.Open(path)
		tools.ErrPanic(err)

		data, err := io.ReadAll(reader)
		tools.ErrPanic(err)

		p := models.PPI{}

		err = json.Unmarshal(data, &p)
		tools.ErrPanic(err)

		for _, pp := range p {
			for _, DevicePpiSamplesList12 := range pp.DevicePpiSamplesList {
				for _, sample := range DevicePpiSamplesList12.PpiSamples {
					sampleTime := time.Time(sample.SampleDateTime)
					influxclient.WritePPIPoint(DevicePpiSamplesList12.DeviceID, sample.PulseLength, sampleTime)
				}
			}
		}
		fmt.Printf("file %s done\n", path)
	}(path)
	return nil
}
