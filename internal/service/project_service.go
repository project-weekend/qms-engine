package service

import (
	"context"

	"github.com/project-weekend/qms-engine/internal/model"
)

type IProjectService interface {
	CreateProject(ctx context.Context, request *model.CreateProjectRequest) (*model.CreateProjectResponse, error)
}
