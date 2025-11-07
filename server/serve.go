package server

import (
	"github.com/sirupsen/logrus"

	"github.com/project-weekend/qms-engine/server/config"
)

func Serve() {
	// Load configuration from service-config.json
	appConfig := config.LoadConfig()
	logger := logrus.New()
	db := config.NewDatabase(appConfig, logger)
	ginEngine := config.NewGinEngine(appConfig, logger)
	if err := config.StartGinServer(ginEngine, appConfig, logger); err != nil {
		logger.Fatalf("Failed to start HTTP server: %v", err)
	}
}
