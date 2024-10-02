package database

import (
	"fmt"
	"polar_reflow/configuration"
	"polar_reflow/database/influxClient"
	"polar_reflow/logger"
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

func InitDB(dbConfig configuration.Database) {
	switch dbConfig.DBType {
	case "influx":
		influxclient.InitInflux(fmt.Sprintf("http://%s", dbConfig.Host), dbConfig.Token, dbConfig.Database, dbConfig.Table)
		Write = func(d models.DBPPI) {
			influxclient.WritePPIPoint(d)
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
		logger.Warning("Mongo is not in operable state. Choose Influx instead")
		return
		//mongoClient.CreateClient(fmt.Sprintf("mongodb://%s:%s@%s", dbConfig.User, dbConfig.Password, dbConfig.Host), dbConfig.Database, dbConfig.Table)
		//Write = func(d models.DBPPI) {
		//	mongoClient.WritePPIPoints(d)
		//}
		//Get = func(t1, t2 time.Time) (result []models.DBPPI) {
		//	cursor := mongoClient.QueryPPI(tools.FormatTime(t1), tools.FormatTime(t2))
		//	err := cursor.All(context.TODO(), &result)
		//	if err != nil {
		//		logger.Error(err.Error())
		//	}
		//	return
		//}
		//Flush = mongoClient.Flush
	default:
		logger.Panic("No database selected")
	}
}

func GetWithTimeAsString(t1, t2 string) []models.DBPPI {
	return Get(tools.ParseTime(t1), tools.ParseTime(t2))
}
