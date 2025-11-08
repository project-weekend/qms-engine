package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/project-weekend/qms-engine/handlers"
	"github.com/project-weekend/qms-engine/internal/repository/mysql"
	"github.com/project-weekend/qms-engine/internal/service/project"
	"github.com/project-weekend/qms-engine/server/config"
	"github.com/sirupsen/logrus"
)

type AppBootstrap struct {
	Config    *config.Config
	Logger    *logrus.Logger
	DB        *sqlx.DB
	Validate  *validator.Validate
	AppEngine *gin.Engine
}

func Bootstrap(app *AppBootstrap) {
	// setup repository
	projectRepository := mysql.NewProjectRepository(app.Logger)

	// setup service
	projectService := project.NewProjectService(app.Logger, app.DB, projectRepository)

	// service injection
	services := handlers.NewQMSEngineService(app.Logger, app.Validate, projectService)

	routeConfig := handlers.RouteConfig{
		AppEngine:        app.AppEngine,
		QMSEngineService: services,
	}

	routeConfig.RegisterRoutes()
}
