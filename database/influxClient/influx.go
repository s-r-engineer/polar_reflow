package influxclient

import (
	"context"
	"fmt"
	"polar_reflow/logger"
	"polar_reflow/models"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

const (
	ppiBucketSuffix   = "_ppi"
	sleepBucketSuffix = "_sleep"

	sleepMeasurement      = "sleep"
	ppiMeasurement        = "ppi"
	baseTemplateWithLimit = `from(bucket:"%s")|> range(start: %sZ, stop: %sZ) |> filter(fn: (r) => r._measurement == "%s" ) |> limit(n: %d, offset: %d)`
)

var (
	client         influxdb2.Client
	writerAPIPPI   api.WriteAPI
	writerAPISleep api.WriteAPI
	queryAPI       api.QueryAPI
	bucketAPI      api.BucketsAPI
	orgAPI         api.OrganizationsAPI

	globalBucket    string
	globalOrg       string
	globalOrgObject *domain.Organization

	ppiBucket   string
	sleepBucket string
)

func InitInflux(influxAddress, token, org, bucket string) {
	var err error
	globalBucket = bucket
	globalOrg = org
	client = influxdb2.NewClientWithOptions(
		influxAddress,
		token,
		influxdb2.DefaultOptions().SetBatchSize(830000))

	for {
		running, err := client.Ping(context.Background())
		if err != nil {
			logger.Error(err.Error())
		}
		if running {
			break
		}
		time.Sleep(time.Second)
	}
	queryAPI = client.QueryAPI(globalOrg)

	bucketAPI = client.BucketsAPI()
	orgAPI = client.OrganizationsAPI()
	globalOrgObject, err = orgAPI.FindOrganizationByName(context.Background(), globalOrg)
	if err != nil {
		logger.Error(err.Error())
	}
	ppiBucket = globalBucket + ppiBucketSuffix
	err = CreateBucket(ppiBucket, false)
	if err != nil {
		logger.Error(err.Error())
	}
	sleepBucket = globalBucket + sleepBucketSuffix
	err = CreateBucket(sleepBucket, false)
	if err != nil {
		logger.Error(err.Error())
	}
	writerAPIPPI = client.WriteAPI(globalOrg, ppiBucket)
	writerAPISleep = client.WriteAPI(globalOrg, sleepBucket)
}

func CreateBucket(name string, force bool) error {
	b, err := bucketAPI.FindBucketByName(context.Background(), name)
	if err == nil && force {
		err = bucketAPI.DeleteBucket(context.Background(), b)
		if err != nil {
			return err
		}
	}
	_, err = bucketAPI.CreateBucket(context.Background(), &domain.Bucket{OrgID: globalOrgObject.Id, Name: name})
	return err
}

func Flush() {
	writerAPIPPI.Flush()
}

func WritePPIPoint(d models.DBPPI) {
	writerAPIPPI.WritePoint(influxdb2.NewPoint(ppiMeasurement,
		map[string]string{},
		map[string]interface{}{"ppi": d.Value},
		d.TimePoint))
}

func WriteHRVPoint(timeTag, method string, data float64, startTime time.Time) {
	writerAPIPPI.WritePoint(influxdb2.NewPoint("hrv",
		map[string]string{"timeScaleInMinutes": timeTag, "method": method},
		map[string]interface{}{"hrv": data},
		startTime))
}

func WriteSleepPoint(s models.SleepResult) {
	writerAPISleep.WritePoint(influxdb2.NewPoint(sleepMeasurement,
		map[string]string{},
		s.ToInflux(),
		time.Time(s.Night)))
}

func QueryPPI(startTime, endTime string, offset int, limit int) *api.QueryTableResult {
	return query(ppiBucket, ppiMeasurement, startTime, endTime, offset, limit)
}

func QuerySleep(startTime, endTime string, offset int, limit int) *api.QueryTableResult {
	return query(sleepBucket, sleepMeasurement, startTime, endTime, offset, limit)
}

func query(bucket, measurement, startTime, endTime string, offset int, limit int) *api.QueryTableResult {
	q := fmt.Sprintf(baseTemplateWithLimit,
		bucket,
		startTime,
		endTime, measurement, limit, offset)
	response, err := queryAPI.Query(context.Background(), q)
	if err != nil {
		logger.Error(err.Error())
	}
	return response
}
