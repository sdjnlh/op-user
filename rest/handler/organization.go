package handler

import (
	"context"

	"code.letsit.cn/go/common"
	"code.letsit.cn/go/common/log"
	"code.letsit.cn/go/common/web"
	"code.letsit.cn/go/op-user/model"
	"code.letsit.cn/go/op-user/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrgAPI struct {
	*web.RestHandler
}

func NewOrgAPI() *OrgAPI {
	return &OrgAPI{
		RestHandler: web.DefaultRestHandler,
	}
}

func (api *OrgAPI) Filter(c *gin.Context) {
	var form = &model.OrgFilter{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	result := &common.FilterResult{Result: common.Result{Data: &[]model.Org{}}}
	err = service.Org.Filter(context.Background(), form, result)
	log.Logger.Debug("list org", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *OrgAPI) Get(c *gin.Context) {
	id := c.Param("id")
	log.Logger.Debug("id>>>>>>", zap.Any("id", id))

	if id == "" {
		api.BadRequest(c)
		return
	}
	var err error
	result := &common.Result{Data: &model.Org{}}
	err = service.Org.Get(context.Background(), &id, result)
	log.Logger.Debug("get role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *OrgAPI) Insert(c *gin.Context) {

	form := &model.Org{}
	file, err := c.FormFile("file")
	form.Name = c.PostForm("name")
	form.Phone = c.PostForm("phone")
	form.Address = c.PostForm("address")
	form.Describe = c.PostForm("describe")
	form.Email = c.PostForm("email")
	form.FullName = c.PostForm("fullName")
	form.Roles = c.PostFormArray("roles")
	log.Logger.Debug("form ", zap.Any("form", form))
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	fileName := file.Filename
	//url :="replace image rul address "
	fileUrl := "/work/" + fileName
	if err = c.SaveUploadedFile(file, fileUrl); err != nil {
		log.Logger.Error("保存失败", zap.Error(err))
		api.BadRequestWithError(c, err)
	}
	form.ImageLogo = fileUrl
	result := &common.Result{Data: &model.Org{}}
	err = service.Org.Insert(context.Background(), form, result)

	log.Logger.Debug("create role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *OrgAPI) Update(c *gin.Context) {

	form := &model.Org{}
	err := api.Bind(c, form)

	if err != nil {
		log.Logger.Error("update role", zap.Any(" err:", err))
		api.BadRequestWithError(c, err)
		return
	}
	log.Logger.Debug("update role", zap.Any(" form:", form))
	result := &common.Result{Data: &model.Org{}}
	err = service.Org.Update(context.Background(), form, result)

	log.Logger.Debug("update role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *OrgAPI) Delete(c *gin.Context) {
	id := c.Param("id")

	var err error
	result := &common.Result{Data: &model.Org{}}

	err = service.Org.Delete(context.Background(), id, result)
	if err != nil {
		log.Logger.Error("delete role", zap.Any(" err:", err))
		api.BadRequestWithError(c, err)
		return
	}
	result.Ok = true
	log.Logger.Debug("delete role", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *OrgAPI) Register(router gin.IRouter) {
	v1 := router.Group("/v1", web.UserInterceptor)
	v1.GET("/orgs", api.Filter)
	v1.GET("/orgs/:id", api.Get)
	v1.POST("/orgs", api.Insert)
	v1.PUT("/orgs/:id", api.Update)
	v1.DELETE("/orgs/:id", api.Delete)
}
