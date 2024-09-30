package mygin

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"polar_reflow/hrv"
	importData "polar_reflow/import"
	"polar_reflow/logger"
	"polar_reflow/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.DefaultWriter = logger.GinWriter{}
    gin.DefaultErrorWriter = logger.GinErrWriter{}
	engine := gin.New()
	engine.Use(logger.LoggerForGin)
	engine.Use(auth)
	engine.PUT("/uploaddata", func(ctx *gin.Context) {
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
		ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded successfully!", file.Filename))
		if err := tools.UnpackArchive(filePath, "/tmp/"); err != nil {
			logger.Error(err.Error())
			ctx.String(http.StatusInternalServerError, "Unable to unpack zip file: %s", err.Error())
			return
		}

		// Optionally, delete the temp file after extraction
		os.Remove(filePath)
		importData.ImportFiles("/tmp")
	})
	engine.GET("/hrv/last5min", getRealHRVRMSSD)
	engine.GET("/hrv/5minforperiod", getRealHRVRMSSDMinByMin)
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	engine.Run(":6969")
}

func auth(ctx *gin.Context) {
	ctx.Next()
}

func getRealHRVRMSSD(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	tools.Dumper(params)
	tempValue := params.Get("from")
	value, err := strconv.Atoi(tempValue)
	logger.Error(err.Error())
	if value == 0 {
		ctx.AbortWithError(512, fmt.Errorf(""))
	}
	ctx.JSON(200, hrv.Get5MinRMSSDFromPoint(value))
}

func getRealHRVRMSSDMinByMin(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	tools.Dumper(params)
	from, err := time.Parse(time.RFC3339, params.Get("from"))
	logger.Error(err.Error())
	to, err := time.Parse(time.RFC3339, params.Get("to"))
	logger.Error(err.Error())
	ctx.JSON(200, hrv.Get5MinRMSSDFromtimeToTime(from, to))
}

type structure1 struct {
	Timestamp int64
	Value     int
}
type structure2 []structure1

func t1Demo(ctx *gin.Context) {
	tools.Dumper(ctx.Request.URL.Query())
	str2 := structure2{}
	for ctime := time.Now().Add(time.Duration(time.Hour * -6)).Unix(); ctime < time.Now().Unix(); ctime++ {
		str2 = append(str2, structure1{Timestamp: ctime, Value: rand.IntN(100)})
	}

	ctx.JSON(200, str2)
}
