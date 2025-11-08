package project

import (
	"context"
	"database/sql"
	"errors"
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
			p.Logger.WithContext(ctx).Error(logTag, "CreateProject GetByName error: ", err)
			return nil, common.NewServiceError(common.ErrCode_InternalServerError, nil)
		}
	} else if existingProject != nil {
		p.Logger.WithContext(ctx).Warn(logTag, "CreateProject: project name already exists: ", request.Name)
		return nil, common.NewServiceError(common.ErrCode_Forbidden, nil)
	}

	project := &entity.Project{
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	savedProject, err := p.ProjectRepository.Save(tx, project)
	if err != nil {
		p.Logger.WithContext(ctx).Error(logTag, "Save project error: ", err)
		return nil, common.NewServiceError(common.ErrCode_InternalServerError, nil)
	}

	err = tx.Commit()
	if err != nil {
		p.Logger.WithContext(ctx).Error(logTag, "Commit project error: ", err)
		return nil, common.NewServiceError(common.ErrCode_InternalServerError, nil)
	}

	return converter.ProjectToResponse(savedProject), nil
}
