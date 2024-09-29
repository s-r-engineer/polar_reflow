package configuration

import (
	"polar_reflow/tools"

	"github.com/kelseyhightower/envconfig"
)

const ENV_PREFIX = "POLAR_REFLOW_"

func Configure() {

}

func ParseConfigFile(path string, currentConfig Config) Config {

}

func ParseEnv(currentConfig Config) Config {
	tools.ErrPanic(envconfig.Process(ENV_PREFIX, &currentConfig))
	return currentConfig
}

func ParseArguments(currentConfig Config) Config {

}

// pathToFindFilesIn         = flag.String("path", os.Getenv("POLAR_REFLOW_PATH"), "")
// serve                     = flag.Bool("serve", "true" == os.Getenv("POLAR_REFLOW_SERVE"), "")
// database                  = flag.String("database", os.Getenv("POLAR_REFLOW_DATABASE"), "Database backend. Could be mongo or influx")
// reinit                    = flag.Bool("reinit", "true" == os.Getenv("POLAR_REFLOW_REINIT"), "")
// dbAddress                 = flag.String("db_address", "http://influx:8086", "")
// excludeRmssd              = flag.Bool("normssd", "true" == os.Getenv("POLAR_REFLOW_NORMSSD"), "")
// excludeSddn               = flag.Bool("nosddn", "true" == os.Getenv("POLAR_REFLOW_NOSDDN"), "")
// token                     = flag.String("token", os.Getenv("POLAR_REFLOW_TOKEN"), "")
// org                       = flag.String("org", os.Getenv("POLAR_REFLOW_ORG"), "")
// bucket                    = flag.String("bucket", os.Getenv("POLAR_REFLOW_BUCKET"), "")
// startTimeString           = flag.String("start", "2020-01-01T00:00:00Z", "")
// finaltime                 = flag.String("end", time.Now().Format("2006-01-02T15:04:05Z"), "")
// parallelismForCalculating = flag.Int("paralel", 16, "")
