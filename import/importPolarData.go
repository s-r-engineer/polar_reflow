package importData

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"polar_reflow/database"
	"polar_reflow/logger"
	"polar_reflow/tools"
	"regexp"

	"polar_reflow/models"
	"polar_reflow/syncronization"
	"time"
)

const (
	ppiRegex   = `^ppi_samples.*\.json$`
	sleepRegex = `^sleep_score.*\.json$`
)

func ImportFiles(pathToLookIn string) {
	logger.Info("Starting reading files")
	absPath, err := filepath.Abs(pathToLookIn)
	if err != nil {
		logger.Error(err.Error())
	}
	aqquire, release := syncronization.CreateSemaphoreInstance(4)
	add, done, wait := syncronization.CreateWGInstance()
	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() {
			return nil
		}

		return importFile(path, aqquire, release, add, done)
	})
	if err != nil {
		logger.Error(err.Error())
	}
	wait()
	logger.Info("flushing data")

	database.Flush()
}

func importFile(path string, aqquire func() error, release, add, done func()) error {
	add()
	go func(path string) {
		defer done()
		ppiFile, err := regexp.MatchString(ppiRegex, filepath.Base(path))
		if err != nil {
			logger.Error(err.Error())
		}
		sleepFile, err := regexp.MatchString(sleepRegex, filepath.Base(path))
		if err != nil {
			logger.Error(err.Error())
		}
		if (!ppiFile && !sleepFile) || (ppiFile && sleepFile) {
			return
		}
		err = aqquire()
		if err != nil {
			logger.Error(err.Error())
			return
		}
		defer release()

		logger.Infof("file %s parsing\n", path)

		reader, err := os.Open(path)
		if err != nil {
			logger.Error(err.Error())
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			logger.Error(err.Error())
		}
		if ppiFile {
			p := models.PPI{}
			err = json.Unmarshal(data, &p)
			if err != nil {
				logger.Error(err.Error())
			}
			for _, pp := range p {
				for _, DevicePpiSamplesList12 := range pp.DevicePpiSamplesList {
					for _, sample := range DevicePpiSamplesList12.PpiSamples {
						sampleTime := time.Time(sample.SampleDateTime)
						database.WritePPI(models.DBPPI{Value: float64(sample.PulseLength), TimePoint: sampleTime})
					}
				}
			}
			database.Flush()
		} else if sleepFile {
			var p models.SleepResults
			err = json.Unmarshal(data, &p)
			if err != nil {
				logger.Error(err.Error())
			}
			for _, sleepResult := range p {
				database.WriteSleep(sleepResult)
			}
			database.Flush()
		}

		logger.Infof("file %s done\n", path)
	}(path)
	return nil
}

func UploadGinHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, "File upload error: %s", err.Error())
		return
	}
	filePath := "/tmp/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.String(http.StatusInternalServerError, "Unable to save the file: %s", err.Error())
		return
	}

	unpackFolder := "/tmp/uploaded"
	if err := tools.UnpackArchive(filePath, unpackFolder); err != nil {
		logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, "Unable to unpack zip file: %s", err.Error())
		return
	}
	ImportFiles("/tmp/uploaded")
	os.Remove(filePath)
	os.Remove(unpackFolder)
	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded successfully processed", file.Filename))
}
