package opu

import (
	"code.letsit.cn/go/common"
	"code.letsit.cn/go/common/app"
	"code.letsit.cn/go/op-user/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

var Service = app.NewService("opu-service")
var Api = app.NewWeb("opu-api")
func UserBuilder(fields map[string]interface{}, c *gin.Context) interface{} {
	user := &model.User{}
	if fields[common.UserIdKey] != nil {
		user.Id, _ = strconv.ParseInt(fields[common.UserIdKey].(string), 10, 64)
	}
	//if fields[common.UserTypeKey] != nil {
	//	user.Type = fields[common.UserTypeKey].(string)
	//}
	if fields[common.UserNicknameKey] != nil {
		user.Nickname = fields[common.UserNicknameKey].(string)
	}
	//if fields[common.UserOrgIdKey] != nil {
	//	user.OrgId, _ = strconv.ParseInt(fields[common.UserOrgIdKey].(string), 10, 64)
	//}

	c.Set(common.UserKey, user)
	return user
}

func User(c *gin.Context) *model.User {
	user, ok := c.Get(common.UserKey)
	if !ok {
		return nil
	}
	return user.(*model.User)
}
