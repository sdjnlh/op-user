package rest

import (
	"code.letsit.cn/go/op-user/rest/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAPIs(router gin.IRouter) {
	handler.NewUserAPI().Register(router)
	handler.NewRoleAPI().Register(router)
	handler.NewGroupAPI().Register(router)
	handler.NewOrgAPI().Register(router)
	handler.NewMiniUserApi().Register(router)
}

func RegisterAPIsWithLogInterceptor(router gin.IRouter) {
	RegisterAPIs(router)
	handler.NewSysLogAPI().Register(router)
}
