package opu

import (
	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/app"
	"github.com/sdjnlh/op-user/model"
	"strconv"
)

var Service = app.NewService("opu-service")
var Api = app.NewWeb("opu-api")

func UserBuilder(fields map[string]interface{}, c *gin.Context) interface{} {
	user := &model.User{}
	if fields[communal.UserIdKey] != nil {
		user.Id, _ = strconv.ParseInt(fields[communal.UserIdKey].(string), 10, 64)
	}
	//if fields[communal.UserTypeKey] != nil {
	//	user.Type = fields[communal.UserTypeKey].(string)
	//}
	if fields[communal.UserNicknameKey] != nil {
		user.Nickname = fields[communal.UserNicknameKey].(string)
	}
	//if fields[communal.UserOrgIdKey] != nil {
	//	user.OrgId, _ = strconv.ParseInt(fields[communal.UserOrgIdKey].(string), 10, 64)
	//}

	c.Set(communal.UserKey, user)
	return user
}

func User(c *gin.Context) *model.User {
	user, ok := c.Get(communal.UserKey)
	if !ok {
		return nil
	}
	return user.(*model.User)
}
