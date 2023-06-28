package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/web"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/service"
	"go.uber.org/zap"
)

type RoleAPI struct {
	*web.RestHandler
}

func NewRoleAPI() *RoleAPI {
	return &RoleAPI{
		RestHandler: web.DefaultRestHandler,
	}
}

func (api *RoleAPI) Filter(c *gin.Context) {
	var form = &model.RoleFilter{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	strOrg, _ := c.Cookie(communal.UserOrgIdKey)
	orgid, _ := strconv.ParseInt(strOrg, 10, 64)
	form.OwnerId = orgid
	result := &communal.FilterResult{Result: communal.Result{Data: &[]model.Role{}}}
	err = service.Role.Filter(context.Background(), form, result)
	log.Logger.Debug("list role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *RoleAPI) Get(c *gin.Context) {
	id := c.Param("id")
	log.Logger.Debug("id>>>>>>", zap.Any("id", id))

	if id == "" {
		api.BadRequest(c)
		return
	}
	var err error
	result := &communal.Result{Data: &model.RoleDTO{}}
	err = service.Role.Get(context.Background(), &id, result)
	log.Logger.Debug("get role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *RoleAPI) Insert(c *gin.Context) {
	form := &model.RoleDTO{}
	err := api.Bind(c, form)
	log.Logger.Debug("form ", zap.Any("form", form))
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	strOrg, _ := c.Cookie(communal.UserOrgIdKey)
	orgid, _ := strconv.ParseInt(strOrg, 10, 64)
	form.OwnerId = orgid
	result := &communal.Result{Data: &model.Role{}}
	err = service.Role.Insert(context.Background(), form, result)

	log.Logger.Debug("create role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *RoleAPI) Update(c *gin.Context) {

	form := &model.RoleDTO{}
	err := api.Bind(c, form)

	if err != nil {
		log.Logger.Error("update role", zap.Any(" err:", err))
		api.BadRequestWithError(c, err)
		return
	}

	log.Logger.Debug("update role", zap.Any(" form:", form))
	result := &communal.Result{Data: &model.Role{}}
	err = service.Role.Update(context.Background(), form, result)

	log.Logger.Debug("update role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *RoleAPI) Delete(c *gin.Context) {
	id := c.Param("id")

	var err error
	result := &communal.Result{Data: &model.Role{}}

	err = service.Role.Delete(context.Background(), id, result)
	if err != nil {
		log.Logger.Error("delete role", zap.Any(" err:", err))
		api.BadRequestWithError(c, err)
		return
	}
	result.Ok = true
	log.Logger.Debug("delete role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *RoleAPI) Register(router gin.IRouter) {
	v1 := router.Group("/v1", web.UserInterceptor)
	v1.GET("/roles", api.Filter)
	v1.GET("/roles/:id", api.Get)
	v1.POST("/roles", api.Insert)
	v1.PUT("/roles/:id", api.Update)
	v1.DELETE("/roles/:id", api.Delete)
}
