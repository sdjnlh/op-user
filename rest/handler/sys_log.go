package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/web"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/service"
)

type sysLogAPI struct {
	*web.RestHandler
}

func NewSysLogAPI() *sysLogAPI {
	return &sysLogAPI{
		web.DefaultRestHandler,
	}
}

func (api *sysLogAPI) List(c *gin.Context) {

	var form = &model.SysLogFilter{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	result := communal.NewFilterResult(&[]model.SysLog{})
	err = service.SysLog.List(form, result)
	if err != nil {
		api.BadRequestWithError(c, err)
	}
	api.ResultWithError(c, result, err)
}

func (api *sysLogAPI) Register(router gin.IRouter) {
	v1 := router.Group("v1/sysLog", web.UserInterceptor)
	v1.GET("", api.List)
}
