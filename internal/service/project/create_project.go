package project

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/project-weekend/qms-engine/internal/common"
	"github.com/project-weekend/qms-engine/internal/entity"
	"github.com/project-weekend/qms-engine/internal/model"
	"github.com/project-weekend/qms-engine/internal/model/converter"
)

func (p *ProjectServiceImpl) CreateProject(ctx context.Context, request *model.CreateProjectRequest) (*model.CreateProjectResponse, error) {
	tx := p.DB.MustBeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	defer tx.Rollback()

	existingProject, err := p.ProjectRepository.GetByName(tx, request.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			p.Logger.ErrorContext(ctx, "CreateProject GetByName error", "tag", logTag, "error", err)
			return nil, common.NewServiceError(common.ErrCode_InternalServerError, nil)
		}
	} else if existingProject != nil {
		p.Logger.WarnContext(ctx, "CreateProject: project name already exists", "tag", logTag, "name", request.Name)
		return nil, common.NewServiceError(common.ErrCode_Forbidden, nil)
	}

	project := &entity.Project{
		Name:        strings.ToLower(request.Name),
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	savedProject, err := p.ProjectRepository.Save(tx, project)
	if err != nil {
		p.Logger.ErrorContext(ctx, "Save project error", "tag", logTag, "error", err)
		return nil, common.NewServiceError(common.ErrCode_InternalServerError, nil)
	}

	err = tx.Commit()
	if err != nil {
		p.Logger.ErrorContext(ctx, "Commit project error", "tag", logTag, "error", err)
		return nil, common.NewServiceError(common.ErrCode_InternalServerError, nil)
	}

	return converter.ProjectToResponse(savedProject), nil
}
