package importData

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	influxclient "polar_reflow/influxClient"
	"polar_reflow/logger"
	"polar_reflow/models"
	"polar_reflow/syncronization"
	"regexp"
	"time"
)

func ImportFiles(pathToLookIn string) {
	logger.Info("Starting reading files")
	absPath, err := filepath.Abs(pathToLookIn)
	logger.Error(err.Error())
	aqquire, release := syncronization.CreateSemaphoreInstance(4)
	add, done, wait := syncronization.CreateWGInstance()
	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		return importFile(path, aqquire, release, add, done)
	},
	)
	logger.Error(err.Error())

	wait()
	logger.Info("flushing data")

	influxclient.Flush()
}

func importFile(path string, aqquire func() error, release, add, done func()) error {
	m, err := regexp.MatchString(`^.*ppi_.*\.json$`, filepath.Base(path))
	logger.Error(err.Error())
	if !m {
		return nil
	}
	add()
	go func(path string) {
		defer done()
		logger.Error(aqquire().Error())

		defer release()

		logger.Infof("file %s parsing\n", path)

		reader, err := os.Open(path)
		logger.Error(err.Error())

		data, err := io.ReadAll(reader)
		logger.Error(err.Error())

		p := models.PPI{}

		err = json.Unmarshal(data, &p)
		logger.Error(err.Error())

		for _, pp := range p {
			for _, DevicePpiSamplesList12 := range pp.DevicePpiSamplesList {
				for _, sample := range DevicePpiSamplesList12.PpiSamples {
					sampleTime := time.Time(sample.SampleDateTime)
					influxclient.WritePPIPoint(DevicePpiSamplesList12.DeviceID, sample.PulseLength, sampleTime)
				}
			}
		}
		logger.Infof("file %s done\n", path)
	}(path)
	return nil
}
