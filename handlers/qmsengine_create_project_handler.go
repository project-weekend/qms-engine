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
		s.Logger.ErrorContext(ctx, "Failed to parse request body", "tag", logTag, "error", err)
		serviceErr := common.NewServiceError(common.ErrCode_BadRequest, nil)
		ctx.AbortWithStatusJSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	if err = s.Validator.Struct(request); err != nil {
		s.Logger.ErrorContext(ctx, "Validation error", "tag", logTag, "error", err)
		serviceErr := common.NewServiceError(common.ErrCode_BadRequest, common.ParseValidationErrors(err))
		ctx.AbortWithStatusJSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	projectResponse, err := s.ProjectService.CreateProject(ctx, request)
	if err != nil {
		s.Logger.ErrorContext(ctx, "CreateProject error", "tag", logTag, "error", err)
		serviceErr := common.AsServiceError(err)
		ctx.AbortWithStatusJSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, projectResponse)
}
