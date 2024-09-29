package database

import (
	"context"
	influxclient "polar_reflow/database/influxClient"
	"polar_reflow/database/mongoClient"
	"polar_reflow/models"
	"polar_reflow/tools"

	"time"
)

var Write func(models.DBPPI)
var Get func(time.Time, time.Time) []models.DBPPI
var Flush func()

type DB interface {
	Write(models.DBPPI)
	Read(time.Time, time.Time) []models.DBPPI
}

func InitDB(dbType, address, tokenIfInflux, dbOrOrg, bucketOrCollection string) {
	switch dbType {
	case "influx":
		influxclient.InitInflux(address, tokenIfInflux, dbOrOrg, bucketOrCollection)
		Write = func(d models.DBPPI) {
			influxclient.WritePPIPoint(d.Value, d.TimePoint)
		}
		Get = func(t1, t2 time.Time) (result []models.DBPPI) {
			cursor := influxclient.QueryPPI(tools.FormatTime(t1), tools.FormatTime(t2))
			for cursor.Next() {
				result = append(result, models.DBPPI{Value: cursor.Record().Value().(float64), TimePoint: cursor.Record().Time()})
			}
			cursor.Close()
			return
		}
		Flush = influxclient.Flush

	case "mongo":
		mongoClient.CreateClient(address, dbOrOrg, bucketOrCollection)
		Write = func(d models.DBPPI) {
			mongoClient.WritePPIPoint(d.Value, d.TimePoint)
		}
		Get = func(t1, t2 time.Time) (result []models.DBPPI) {
			cursor := mongoClient.QueryPPI(tools.FormatTime(t1), tools.FormatTime(t2))
			tools.ErrPanic(cursor.All(context.TODO(), &result))
			return
		}
		Flush = func() {}
	default:
		return
	}
}

func GetWithTimeAsString(t1, t2 string) []models.DBPPI {
	return Get(tools.ParseTime(t1), tools.ParseTime(t2))
}
