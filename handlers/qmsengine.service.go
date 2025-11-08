package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/project-weekend/qms-engine/internal/service/project"
)

const logTag = "handlers"

// QMSEngineService holds all dependencies for QMS Engine handlers
type QMSEngineService struct {
	Logger         *logrus.Logger
	Validator      *validator.Validate
	ProjectService *project.ProjectServiceImpl
}

func NewQMSEngineService(logger *logrus.Logger, validator *validator.Validate, projectService *project.ProjectServiceImpl) *QMSEngineService {
	return &QMSEngineService{
		Logger:         logger,
		Validator:      validator,
		ProjectService: projectService,
	}
}
