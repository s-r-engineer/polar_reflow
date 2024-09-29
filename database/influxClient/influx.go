package influxclient

import (
	"context"
	"fmt"
	"polar_reflow/logger"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

var (
	client    influxdb2.Client
	writerAPI api.WriteAPI
	queryAPI  api.QueryAPI

	globalBucket string
	globalOrg    string
)

func InitInflux(influxAddress, token, org, bucket string) {
	globalBucket = bucket
	globalOrg = org
	client = influxdb2.NewClientWithOptions(
		influxAddress,
		token,
		influxdb2.DefaultOptions().SetBatchSize(830000))
	writerAPI = client.WriteAPI(globalOrg, globalBucket)
	queryAPI = client.QueryAPI(globalOrg)
}

func ReinitBucket() {
	logger.Info("reinit starting")
	bucketAPI := client.BucketsAPI()
	orgAPI := client.OrganizationsAPI()
	o, err := orgAPI.FindOrganizationByName(context.Background(), globalOrg)
	logger.Error(err.Error())

	b, err := bucketAPI.FindBucketByName(context.Background(), globalBucket)
	if err == nil {
		err = bucketAPI.DeleteBucket(context.Background(), b)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	_, err = bucketAPI.CreateBucket(context.Background(), &domain.Bucket{OrgID: o.Id, Name: globalBucket})
	logger.Error(err.Error())
	logger.Info("reinit done")
}

func Flush() {
	writerAPI.Flush()
}

func WritePPIPoint(pulseLength float64, sampleTime time.Time) {
	writerAPI.WritePoint(influxdb2.NewPoint("ppi",
		map[string]string{},
		map[string]interface{}{"ppi": pulseLength},
		sampleTime))
}

func WriteHRVPoint(timeTag, method string, data float64, startTime time.Time) {
	writerAPI.WritePoint(influxdb2.NewPoint("hrv",
		map[string]string{"timeScaleInMinutes": timeTag, "method": method},
		map[string]interface{}{"hrv": data},
		startTime))
}

func QueryPPI(startTime, endTime string) *api.QueryTableResult {
	q := fmt.Sprintf(`from(bucket:"%s")|> range(start: %sZ, stop: %sZ) |> filter(fn: (r) => r._measurement == "ppi")`,
		globalBucket,
		startTime,
		endTime)
	response, err := queryAPI.Query(context.Background(), q)
	logger.Error(err.Error())
	return response
}
