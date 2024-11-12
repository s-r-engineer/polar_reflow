package configuration

import (
	"github.com/kelseyhightower/envconfig"
	"polar_reflow/logger"
)

const EnvPrefix = "POLAR_REFLOW_"

func Configure() (c Config) {
	c.Database.DBType = "influx"
	c.Database.Host = "localhost:8086"
	c.Database.Database = "polar_reflow"
	c.Database.Table = "hrv"
	c.Database.User = "polar_reflow"
	c.Database.Password = "polar_reflow"
	c.Database.Token = "sdwefrejktkyukytjrjtrw"
	c.API = API{BindAddress: "0.0.0.0:6969"}
	c.Engine = Engine{Parallel: 16}
	return parseEnv(c)
}

func parseEnv(currentConfig Config) Config {
	err := envconfig.Process(EnvPrefix, &currentConfig)
	if err != nil {
		logger.Panic(err.Error())
	}
	return currentConfig
}
