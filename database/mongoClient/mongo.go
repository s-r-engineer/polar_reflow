package mongoClient

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"polar_reflow/logger"
	"polar_reflow/models"
	"polar_reflow/syncronization"
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

	inserMany         = 100000
	insertManyStorage []interface{}

	lock, unlock = syncronization.CreateMutexInstance()

	fALSE = false
)

func CreateClient(uri, database, coll string) {
	globalDatabase = database
	globalCollection = coll
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error(err.Error())
	}
	//granularity := "seconds"
	err = client.Database(globalDatabase).CreateCollection(context.TODO(), globalCollection)
	//err = client.Database(globalDatabase).CreateCollection(context.TODO(), globalCollection, &options.CreateCollectionOptions{TimeSeriesOptions: &options.TimeSeriesOptions{
	//	TimeField:   "timePoint",
	//	Granularity: &granularity,
	//}})
	if err != nil {
		logger.Error(err.Error())
	}

	collection = client.Database(globalDatabase).Collection(globalCollection)
	//indexModel := mongo.IndexModel{
	//	Keys: bson.M{
	//		"timePoint": 1,
	//	},
	//	Options: options.Index().SetUnique(true), // Make the index unique
	//}
	//_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	//if err != nil {
	//	logger.Panic(err.Error())
	//}

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
	lock()
	defer unlock()
	_, err := collection.InsertOne(context.TODO(), models.MongoDBPPI{ID: genID(ppi, sampleTime), Value: ppi, TimePoint: sampleTime})
	if err != nil {
		logger.Error(err.Error())
	}
}

func genID(ppi float64, sampleTime time.Time) string {
	hash := sha3.New256()
	hash.Write([]byte(fmt.Sprintf("%d-%d", ppi, sampleTime.UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func WritePPIPoints(d models.DBPPI) {
	lock()
	defer unlock()
	insertManyStorage = append(insertManyStorage, models.MongoDBPPI{ID: genID(d.Value, d.TimePoint), Value: d.Value, TimePoint: d.TimePoint})
	if len(insertManyStorage) >= inserMany {
		nonBlockingFlush()
	}
}

func nonBlockingFlush() {
	if len(insertManyStorage) > 0 {
		collection.InsertMany(context.TODO(), insertManyStorage, &options.InsertManyOptions{Ordered: &fALSE})

		//if err != nil && !isDup(err) {
		//	logger.Error(err.Error())
		//	return
		//}
		insertManyStorage = []interface{}{}
	}
}

func Flush() {
	lock()
	defer unlock()
	nonBlockingFlush()
}

func isDup(err error) bool {
	var merr mongo.BulkWriteException
	errors.As(err, &merr)
	for _, we := range merr.WriteErrors {
		if we.Code != 11000 {
			return false
		}
	}
	return true
}
