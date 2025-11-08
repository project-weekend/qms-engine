package handlers

import "github.com/gin-gonic/gin"

type RouteConfig struct {
	AppEngine *gin.Engine
	*QMSEngineService
}

func (r *RouteConfig) RegisterRoutes() {
	r.AppEngine.POST("/v1/project", r.CreateProject)
}
