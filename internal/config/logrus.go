package config

import (
	"github.com/project-weekend/qms-engine/server/config"
	"github.com/sirupsen/logrus"
)

func NewLogger(appCfg *config.Config) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(appCfg.Logger.LogLevel))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
