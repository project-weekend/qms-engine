package server

import (
	"fmt"
	"log"

	"github.com/project-weekend/qms-engine/internal/config"
)

func Serve() {
	appConfig := config.LoadConfig()
	logger := config.NewLogger(appConfig)
	db := config.NewDatabase(appConfig, logger)
	validator := config.NewValidator()
	appEngine := config.NewGinEngine(appConfig, logger)

	config.Bootstrap(&config.AppBootstrap{
		Config:    appConfig,
		Logger:    logger,
		DB:        db,
		Validate:  validator,
		AppEngine: appEngine,
	})

	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))

	if err := appEngine.Run(addr); err != nil {
		log.Fatal(fmt.Errorf("failed to start http server: %w", err))
	}
}
