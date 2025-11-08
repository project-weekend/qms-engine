package project

import (
	"github.com/jmoiron/sqlx"
	"github.com/project-weekend/qms-engine/internal/repository/mysql"
	"github.com/sirupsen/logrus"
)

const (
	logTag = "service.project"
)

type ProjectServiceImpl struct {
	Logger            *logrus.Logger
	DB                *sqlx.DB
	ProjectRepository *mysql.ProjectRepository
}

func NewProjectService(logger *logrus.Logger, db *sqlx.DB, projectRepository *mysql.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		Logger:            logger,
		DB:                db,
		ProjectRepository: projectRepository,
	}
}
