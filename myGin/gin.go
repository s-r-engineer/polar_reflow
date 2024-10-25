package mygin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polar_reflow/hrv"
	importData "polar_reflow/import"
	"polar_reflow/logger"
	"polar_reflow/sleep"
)

func Run() {
	gin.DefaultWriter = logger.GinWriter{}
	gin.DefaultErrorWriter = logger.GinErrWriter{}
	engine := gin.New()
	engine.Use(logger.LoggerForGin)
	engine.Use(auth)
	engine.PUT("/uploaddata", importData.UploadGinHandler)
	engine.GET("/hrv/5minforperiod", hrv.GetRealHRVRMSSDMinByMin)
	engine.GET("/sleep/forperiod", sleep.GetSleepForPeriod)
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	engine.Run(":6969")
}

func auth(ctx *gin.Context) {
	ctx.Next()
}
