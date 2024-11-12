package influxclient

import (
	"context"
	"fmt"
	"polar_reflow/logger"
	"polar_reflow/models"
	"polar_reflow/tools"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	//"github.com/influxdata/influxdb-client-go/v2/domain"
)

var (
	client      influxdb2.Client
	writerAPI   api.WriteAPI
	writerHRAPI api.WriteAPI
	queryAPI    api.QueryAPI

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
	for {
		ok, err := client.Ping(context.TODO())
		if ok {
			tools.Dumper(ok)
			break
		}
		tools.Dumper(err)
	}
	//ReinitBucket()
	writerAPI = client.WriteAPI(globalOrg, globalBucket)
	writerHRAPI = client.WriteAPI(globalOrg, "HR")
	queryAPI = client.QueryAPI(globalOrg)
	//tools.Dumper(12)
}

func ReinitBucket() {
	logger.Info("reinit starting")
	bucketAPI := client.BucketsAPI()
	orgAPI := client.OrganizationsAPI()
	_, err := orgAPI.FindOrganizationByName(context.Background(), globalOrg)
	//tools.Dumper(orgAPI)
	if err != nil {
		logger.Error(err.Error())
	}

	_, err = bucketAPI.FindBucketByName(context.Background(), "HR")
	//if err != nil {
	//	_, err = bucketAPI.CreateBucket(context.Background(), &domain.Bucket{OrgID: o.Id, Name: "HR"})
	//	if err != nil {
	//		logger.Error(err.Error())
	//	}
	//
	//	//logger.Error(err.Error())
	//	logger.Info("reinit done")
	//}

}

func Flush() {
	writerAPI.Flush()
}

func WriteBTHR(d models.BTHR) {
	//tools.Dumper(writerHRAPI)
	writerHRAPI.WritePoint(influxdb2.NewPoint("bthr",
		map[string]string{"session": d.Session.String()},
		map[string]interface{}{"bthr": d.Value},
		d.TimePoint))
}

func WritePPIPoint(d models.DBPPI) {
	writerAPI.WritePoint(influxdb2.NewPoint("ppi",
		map[string]string{},
		map[string]interface{}{"ppi": d.Value},
		d.TimePoint))
}

func WriteHRVPoint(timeTag, method string, data float64, startTime time.Time) {
	writerAPI.WritePoint(influxdb2.NewPoint("hrv",
		map[string]string{"timeScaleInMinutes": timeTag, "method": method},
		map[string]interface{}{"hrv": data},
		startTime))
}

func QueryPPI(startTime, endTime string, offset int, limit int) *api.QueryTableResult {
	q := fmt.Sprintf(`from(bucket:"%s")|> range(start: %sZ, stop: %sZ) |> filter(fn: (r) => r._measurement == "ppi" ) |> limit(n: %d, offset: %d)`,
		globalBucket,
		startTime,
		endTime, limit, offset)
	response, err := queryAPI.Query(context.Background(), q)
	if err != nil {
		logger.Error(err.Error())
	}
	return response
}
