package sleep

import (
	"polar_reflow/database"
	"polar_reflow/logger"
	"polar_reflow/models"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func GetSleepForPeriod(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	from, err := time.Parse(time.RFC3339, params.Get("from"))
	if err != nil {
		logger.Error(err.Error())
	}
	to, err := time.Parse(time.RFC3339, params.Get("to"))
	if err != nil {
		logger.Error(err.Error())
	}
	scoreType := params.Get("type")
	ctx.JSON(200, getSleepWithParameters(from, to, strings.Split(scoreType, ",")...))
}

func getSleepWithParameters(from, to time.Time, opts ...string) models.SleepResults {
	spew.Dump(opts)
	return database.GetSleep(from, to)
}
