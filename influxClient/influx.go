package influxclient

import (
	"context"
	"fmt"
	"polar_reflow/tools"
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
	fmt.Println("reinit starting")
	bucketAPI := client.BucketsAPI()
	orgAPI := client.OrganizationsAPI()
	o, e := orgAPI.FindOrganizationByName(context.Background(), globalOrg)
	tools.ErrPanic(e)

	b, e := bucketAPI.FindBucketByName(context.Background(), globalBucket)
	if e == nil {
		tools.ErrPanic(bucketAPI.DeleteBucket(context.Background(), b))
	}

	_, e = bucketAPI.CreateBucket(context.Background(), &domain.Bucket{OrgID: o.Id, Name: globalBucket})
	tools.ErrPanic(e)

	fmt.Println("reinit done")
}

func Flush() {
	writerAPI.Flush()
}

func WritePPIPoint(deviceID string, pulseLength int, sampleTime time.Time) {
	writerAPI.WritePoint(influxdb2.NewPoint("ppi",
		map[string]string{"device": deviceID},
		map[string]interface{}{"ppi": float64(pulseLength)},
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
	tools.ErrPanic(err)
	return response
}
