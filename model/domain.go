package model

import (
	"time"

	"code.letsit.cn/go/common"
	"code.letsit.cn/go/common/db"
	"github.com/go-xorm/builder"
	"xorm.io/xorm"
)

type User struct {
	common.DBase `xorm:"extends"`
	Username     string             `json:"username" form:"username"`
	Password     string             `json:"password" form:"password"`
	Nickname     string             `json:"nickname" form:"nickname"`
	Email        string             `json:"email" form:"email"`
	Mcc          string             `json:"-" form:"mcc"`
	Mobile       string             `json:"mobile" form:"mobile"`
	AvatarId     string             `json:"-"`
	Language     int                `json:"-"`
	RoleIds      common.StringArray `json:"roleIds"`
	Groups       common.StringArray
	OrgId        int64 `json:",string"`
	Ext          common.JsonMap
	Status       int
	Path         common.Int64Array //机构路径
}
type Org struct {
	common.DBase `xorm:"extends"`
	Name         string `json:"name",from:"name"`
	FullName     string `json:"fullName",from:"fullName"`
	Address      string `json:"address",from:"address"`
	ImageLogo    string `json:"imageLogo",from:"imageLogo"`
	Phone        string `json:"phone",from:"phone"`
	Describe     string `json:"describe",from:"describe"`
	Email        string `json:"email",from:"email"`
	License      string `json:"license",from:"license"`
	Roles        common.StringArray
	Code         string
}

type JnOrg struct {
	common.DBase `xorm:"extends"`
	Name         string           `json:"name"`
	Area         string           `json:"area" form:"area"`
	Address      string           `json:"address" form:"address"`
	Scale        string           `json:"scale" form:"scale"`
	Mobile       string           `json:"mobile" form:"mobile"`
	ParentId     int64            `json:"parentId,string" form:"parentId,string"`
	ManagerUnit  bool             `form:"managerUnit"`
	Progress     []common.JsonMap `json:"progress" form:"progress"`
}

func (user *User) TableName() string {
	return "op_user"
}

type Role struct {
	Id          string
	Name        string `json:"name" form:"name"`
	Code        string `form:"code"`
	Permissions common.StringArray
	Description string    `json:"description" form:"description"`
	Crt         time.Time `json:"crt"`
	Lut         time.Time `json:"-"`
	Dtd         bool
	OwnerId     int64
	Path        string
}

func (role *Role) TableName() string {
	return "op_role"
}

type RoleFilter struct {
	Keyword string `json:"k" form:"k"`
	common.Page
	OwnerId int64
}

type Permission struct {
	Id          string `form:"id"`
	Name        string `json:"name"`
	ParentId    string `form:"parentId"`
	Path        common.StringArray
	Description string    `json:"description"`
	Crt         time.Time `json:"crt"`
	Lut         time.Time `json:"-"`
	Dtd         bool
}

func (right *Permission) TableName() string {
	return "op_permission"
}

type PermissionFilter struct {
	Keyword string `json:"k" form:"k"`
	common.Page
}

func (filter *PermissionFilter) Apply(session *xorm.Session) {
	session.Where(builder.Eq{"status": db.STATUS_COMMON_OK})
	if filter.Keyword != "" {
		session.And("name like ?", "%"+filter.Keyword+"%")
	}
	session.Limit(filter.Limit(), filter.Skip())
}

type Group struct {
	Id          string `form:"id"`
	Name        string `json:"name"`
	Tags        common.StringArray
	Description string    `json:"description"`
	Crt         time.Time `json:"crt"`
	Lut         time.Time `json:"-"`
	Dtd         bool
}

func (domain *Group) TableName() string {
	return "op_group"
}

type GroupFilter struct {
	Keyword string `json:"k" form:"k"`
	common.Page
}

func (filter *GroupFilter) Apply(session *xorm.Session) {
	session.Where(builder.Eq{"status": db.STATUS_COMMON_OK})
	if filter.Keyword != "" {
		session.And("name like ?", "%"+filter.Keyword+"%")
	}
	session.Limit(filter.Limit(), filter.Skip())
}

type SysLog struct {
	common.DBase `xorm:"extends"`
	Ip           string
	Uri          string
	Uid          int64 `json:",string"`
	Params       string
	Ext          common.JsonMap
	Username     string `form:"username"`
	Method       string
	Agent        string
}

type SysLogFilter struct {
	SysLog      `xorm:"extends"`
	common.Page `xorm:"-"`
	StartTime   time.Time `xorm:"-" form:"startTime"`
	EndTime     time.Time `xorm:"-" form:"endTime"`
}

func (domain *SysLog) TableName() string {
	return "sys_log"
}

type DecodePhone struct {
	EncryptedData string `json:"encryptedData" form:"encryptedData"`
	SessionKey    string `json:"sessionKey" form:"sessionKey"`
	Iv            string `json:"iv" form:"iv"`
	OpenId        string `json:"openId" form:"openId"`
}

type Phone struct {
	PhoneNumber     string `json:"phoneNumber" form:"phoneNumber"`
	PurePhoneNumber string
	CountryCode     string
}
type MiniUser struct {
	common.DBase  `xorm:"extends"`
	Username      string            `json:"username" form:"username"`
	Password      string            `form:"password" json:"password"`
	RoleIds       common.Int64Array `json:"roleIds"` // 角色列表
	Email         string            `json:"email" form:"email"`
	Mobile        string            `json:"mobile" form:"mobile"`
	Status        string            `json:"status" form:"status"`
	Openid        string            `json:"openid" form:"openid"`           // open id
	SessionKey    string            `json:"session_key" form:"session_key"` // session key
	Token         string            `json:"token" form:"token"`             // token access token
	NickName      string            `json:"nickName" form:"nickName"`       // 昵称
	AvatarUrl     string            `json:"avatarUrl"`                      // 头像
	Code          string            `json:"code" form:"code"`
	Province      string            `json:"province" form:"province"`
	Gender        int               `json:"gender" form:"gender"`
	City          string            `json:"city" form:"city"`
	Region        string            `json:"region" form:"region"`
	Url           string            `json:"url"`
	Level         string            `json:"level"`
	Cert          bool              `json:"cert" form:"cert"`         // 是否认证 模式是false
	RoleType      int               `json:"roleType" form:"roleType"` // 用户类型
	Company       string            `json:"company" form:"company"`   //用户企业
	BdId          int64             `json:"bdId,string" form:"bdId"`  // 推荐小程序用户的id
	EncryptedData string            `json:"encryptedData" xorm:"-"`
	ErrMsg        string            `json:"errMsg"  xorm:"-"`
	Iv            string            `json:"iv"  xorm:"-"`
	RowData       string            `json:"rowData"  xorm:"-"`
	Signature     string            `json:"signature"  xorm:"-"`
}
