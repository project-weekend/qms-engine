package project

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/project-weekend/qms-engine/internal/repository/mysql"
)

const (
	logTag = "service.project"
)

type ProjectServiceImpl struct {
	Logger            *slog.Logger
	DB                *sqlx.DB
	ProjectRepository *mysql.ProjectRepository
}

func NewProjectService(logger *slog.Logger, db *sqlx.DB, projectRepository *mysql.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		Logger:            logger,
		DB:                db,
		ProjectRepository: projectRepository,
	}
}
