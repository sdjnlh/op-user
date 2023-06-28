package module

import (
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/op-user/opu"
)

var User = communal.NewModule("User", "op_user", "/v1/users")
var Role = communal.NewModule("Role", "op_role", "/v1/roles")
var Permission = communal.NewModule("Permission", "op_permission", "/v1/permissions")
var Group = communal.NewModule("Group", "op_group", "/v1/groups")
var Org = communal.NewModule("Org", "organization", "/v1/orgs")
var SysLog = communal.NewModule("SysLog", "sys_log", "v1/sysLog")
var MiniUser = communal.NewModule("MiniUser", "mini_user", "")

func init() {
	opu.Service.Register(SysLog, User, Role, Permission, Group, Org, MiniUser)
}
