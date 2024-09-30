package mongoClient

import (
	"context"
	"polar_reflow/logger"
	"polar_reflow/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection

	globalDatabase   string
	globalCollection string
)

func CreateClient(uri, database, coll string) {
	globalDatabase = database
	globalCollection = coll
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error(err.Error())
	}
	granularity := "second"
	err = client.Database(globalDatabase).CreateCollection(context.TODO(), globalCollection, &options.CreateCollectionOptions{TimeSeriesOptions: &options.TimeSeriesOptions{
		TimeField:   "timePoint",
		Granularity: &granularity,
	}})
	if err != nil {
		logger.Error(err.Error())
	}
	collection = client.Database(globalDatabase).Collection(globalCollection)
}

func QueryPPI(startTime, endTime string) *mongo.Cursor {
	filter := bson.M{
		"timePoint": bson.M{
			"$gte": startTime,
			"$lt":  endTime,
		},
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		logger.Error(err.Error())
	}
	return cursor
}

func WritePPIPoint(ppi float64, sampleTime time.Time) {
	_, err := collection.InsertOne(context.TODO(), models.DBPPI{Value: ppi, TimePoint: sampleTime})
	if err != nil {
		logger.Error(err.Error())
	}
}
