package main

import (
	"os"
	"polar_reflow/configuration"
	"polar_reflow/database"
	"polar_reflow/logger"
	mygin "polar_reflow/myGin"
	"polar_reflow/testBT"
)

func init() {
	logger.InitLogger(os.Getenv("POLAR_REFLOW_DEPLOYMENT"))
}

func main() {
	logger.Info("Starting")

	defer logger.Info("Quitting")
	config := configuration.Configure()
	database.InitDB(config.Database)
	testBT.RunBT()
	mygin.Run()
}
