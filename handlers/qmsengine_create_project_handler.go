package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/project-weekend/qms-engine/internal/common"
	"github.com/project-weekend/qms-engine/internal/model"
)

// CreateProject handles project creation
func (s *QMSEngineService) CreateProject(ctx *gin.Context) {
	request := new(model.CreateProjectRequest)
	err := ctx.ShouldBind(request)
	if err != nil {
		s.Logger.WithContext(ctx).Error(logTag, "failed to parse request body")
		serviceErr := common.NewServiceError(common.ErrCode_BadRequest, nil)
		ctx.JSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	if err = s.Validator.Struct(request); err != nil {
		s.Logger.WithContext(ctx).Error(logTag, "Validation error, err", err)
		serviceErr := common.NewServiceError(common.ErrCode_BadRequest, common.ParseValidationErrors(err))
		ctx.AbortWithStatusJSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	projectResponse, err := s.ProjectService.CreateProject(ctx, request)
	if err != nil {
		s.Logger.WithContext(ctx).Error(logTag, "CreateProject error", err)
		serviceErr := common.NewServiceError(common.ErrCode_InternalServerError, nil)
		ctx.AbortWithStatusJSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, projectResponse)
}
