package handler

import (
	"code.letsit.cn/go/common"
	"code.letsit.cn/go/common/web"
	"code.letsit.cn/go/op-user/model"
	"code.letsit.cn/go/op-user/service"
	"github.com/gin-gonic/gin"
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

	result := common.NewFilterResult(&[]model.SysLog{})
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
