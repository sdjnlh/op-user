package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/rpc"
	"github.com/sdjnlh/communal/web"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/opu"
	"github.com/sdjnlh/op-user/service"
	"go.uber.org/zap"
)

type GroupAPI struct {
	*web.RestHandler
}

func NewGroupAPI() *GroupAPI {
	return &GroupAPI{
		RestHandler: web.DefaultRestHandler,
	}
}

func (api *GroupAPI) Filter(c *gin.Context) {
	var form = &model.GroupFilter{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	result := &communal.FilterResult{Result: communal.Result{Data: &[]model.Group{}}}

	if opu.Api.Rpc {
		err = rpc.Call(context.Background(), "Group", "Filter", form, result)
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
	} else {
		err = service.Group.List(context.Background(), form, result)
	}
	if err != nil {
		api.BadRequestWithError(c, err)
	}
	result.Ok = true
	log.Logger.Debug("list group", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *GroupAPI) Get(c *gin.Context) {
	id, err := api.ValidateInt64Id(c)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	result := &communal.Result{Data: &model.Group{}}
	if opu.Api.Rpc {
		err = rpc.Call(context.Background(), "Group", "Get", id, result)
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
	} else {
		err = service.Group.Get(context.Background(), &id, result)
	}

	log.Logger.Debug("get group", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *GroupAPI) Insert(c *gin.Context) {
	form := &model.Group{}
	err := api.Bind(c, form)
	log.Logger.Debug("form ", zap.Any("form", form))
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	result := &communal.Result{Data: &model.Group{}}
	if opu.Api.Rpc {
		err = rpc.Call(context.Background(), "Group", "Insert", form, result)
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
	} else {
		err = service.Group.Create(context.Background(), form, result)
	}

	log.Logger.Debug("create group", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *GroupAPI) Update(c *gin.Context) {
	form := &model.Group{}
	err := api.Bind(c, form)

	if err != nil {
		log.Logger.Error("update group", zap.Any(" err:", err))
		api.BadRequestWithError(c, err)
		return
	}
	log.Logger.Error("update group", zap.Any(" form:", form))
	result := &communal.Result{Data: &model.Group{}}
	if opu.Api.Rpc {
		err = rpc.Call(context.Background(), "Group", "Update", form, result)
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
	} else {
		err = service.Group.Update(context.Background(), form, result)
	}

	log.Logger.Debug("update group", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *GroupAPI) Delete(c *gin.Context) {
	id, err := api.ValidateInt64Id(c)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	result := &communal.Result{Data: &model.Group{}}
	if opu.Api.Rpc {
		err = rpc.Call(context.Background(), "Group", "Delete", id, result)
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
	} else {
		err = service.Group.Delete(context.Background(), &id, result)
	}

	log.Logger.Debug("delete group", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *GroupAPI) Register(router gin.IRouter) {
	v1 := router.Group("/v1", web.UserInterceptor)
	v1.GET("/groups", api.Filter)
	v1.GET("/groups/:id", api.Get)
	v1.POST("/groups", api.Insert)
	v1.PUT("/groups", api.Update)
	v1.DELETE("/groups/:id", api.Delete)
}
