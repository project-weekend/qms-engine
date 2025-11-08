package converter

import (
	"github.com/project-weekend/qms-engine/internal/entity"
	"github.com/project-weekend/qms-engine/internal/model"
)

func ProjectToResponse(entity *entity.Project) *model.CreateProjectResponse {
	return &model.CreateProjectResponse{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
