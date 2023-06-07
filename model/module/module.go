package module

import (
	"code.letsit.cn/go/common"
	"code.letsit.cn/go/op-user/opu"
)

var User = common.NewModule("User", "op_user", "/v1/users")
var Role = common.NewModule("Role", "op_role", "/v1/roles")
var Permission = common.NewModule("Permission", "op_permission", "/v1/permissions")
var Group = common.NewModule("Group", "op_group", "/v1/groups")
var Org = common.NewModule("Org", "organization", "/v1/orgs")
var SysLog = common.NewModule("SysLog", "sys_log", "v1/sysLog")
var MiniUser = common.NewModule("MiniUser", "mini_user", "")

func init() {
	opu.Service.Register(SysLog, User, Role, Permission, Group, Org, MiniUser)
}
