package model

import (
	"github.com/sdjnlh/communal"
	"strconv"
)

const (
	USERNAME_TYPE_MOBILE = iota
	USERNAME_TYPE_EMAIL

	MCC_DEFAULT = "86"
)

type UserDTO struct {
	User        `xorm:"extends"`
	Permissions *communal.StringArray `xorm:"-" json:"permissions,omitempty"`
	Groups      *communal.StringArray `xorm:"-" json:"groups,omitempty"`
	Roles       *[]Role               `xorm:"-"`
	Org         Org                   `xorm:"-"`
	DefaultPath string                `xorm:"-"`
	OrgName     string                `xorm:"-" form:"orgName" json:"orgName"`
}

func (user *UserDTO) State() map[string]interface{} {
	return map[string]interface{}{
		communal.UserIdKey:       strconv.FormatInt(user.Id, 10),
		communal.UserNicknameKey: user.Nickname,
	}
}

type UserFilter struct {
	Keyword  string   `json:"k" form:"k"`
	Roles    []string `form:"roles[]"`
	Groups   []string `form:"groups"`
	OwnerId  int64
	OrgId    int64  `form:"orgId,string" json:"orgId,string"`
	Scope    string `form:"scope"`
	ExistOrg bool   `form:"existOrg"`
	communal.Page
}

type RoleDTO struct {
	Role `xorm:"extends"`
	//RightIds communal.StringArray `xorm:"-"`
}

type RestPassword struct {
	Id          int64  `json:"id,string"`
	OldPassword string `json:"oldpassword" form:"oldpassword"`
	NewPassword string `json:"newpassword" form:"newpassword"`
}

type UserBase struct {
	communal.DBase `xorm:"extends"`
	OrgId          int64  `json:"orgId,string"`
	Username       string `json:"username" form:"username"`
	Password       string `json:"password" form:"password"`
	Nickname       string `json:"nickname" form:"nickname"`
	Email          string `json:"email" form:"email"`
	Mcc            string `json:"-" form:"mcc"`
	Mobile         string `json:"mobile" form:"mobile"`
	Type           string `json:"type" form:"type"`
}
type UUser struct {
	UserBase `xorm:"extends"`
	Roles    communal.StringArray
	AvatarId string `json:"-"`
	Language int    `json:"-"`
	Source   int    `json:"source"`
	Ext      communal.JsonMap
}

func (user *UUser) TableName() string {
	return "u_user"
}

type URole struct {
	communal.DBase `xorm:"extends"`
	OrgId          int64  `form:"orgId,string" json:"orgId,string"`
	Name           string `json:"name" form:"name"`
	Permissions    communal.StringArray
	Description    string `json:"description" form:"description"`
}

func (role *URole) TableName() string {
	return "u_role"
}

type URoleFilter struct {
	Keyword string `json:"k" form:"k"`
	communal.Page
}
type URoleDTO struct {
	URole `xorm:"extends"`
	//RightIds communal.StringArray `xorm:"-"`
}

type UUserFilter struct {
	Keyword string `json:"k" form:"k"`
	communal.Page
}

type UUserDTO struct {
	UUser   `xorm:"extends"`
	Rights  *communal.StringArray `xorm:"-" json:"rights,omitempty"`
	Groups  *communal.StringArray `xorm:"-" json:"groups,omitempty"`
	URoles  []URole               `xorm:"-"`
	OrgName string                `json:"orgName" form:"orgName" xorm:"-"`
}

type ForgetForm struct {
	Email    string
	Code     string
	Language string
}
type OrgFilter struct {
	FullName string
	Name     string
	communal.Page
}
type OrgDTO struct {
	Org `xrom:"extends"`
}

type PhoneCode struct {
	PhoneNumber  string
	Code         string
	Tag          string
	UserId       int64  `json:"userId,string" form:"userId,string"`
	Password     string `json:"password"`
	BdId         int64  `json:"bdId,string" form:"bdId,string"`
	Username     string `json:"username" form:"username"`
	MobileFilter string `json:"mobileFilter" form:"mobileFilter"`
}

const CaptchaTimeout = 600
