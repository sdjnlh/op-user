package service

import (
	"context"
	"encoding/json"
	be "errors"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/errors"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/password"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/model/module"
	"go.uber.org/zap"
)

type UserService struct {
	*communal.Module
}

var User = &UserService{
	Module: module.User,
}

func (service *UserService) Login(ctx context.Context, form *model.User, result *communal.Result) error {
	log.Logger.Debug("user login")
	user := result.Data.(*model.UserDTO)

	username := strings.TrimSpace(form.Username)
	pas := strings.TrimSpace(form.Password)

	has, err := service.Db.Where(" dtd = false").And("mobile = ? or email ilike ? or username ilike ?", username, username, username).Get(user)
	log.Logger.Debug("select user by mobile", zap.Bool("has", has))
	if err != nil {
		log.Logger.Error("failed to get user by mobile", zap.Any("err", err))
		return err
	}

	if user.Id <= 0 {
		result.Failure(&errors.SimpleBizError{
			Code: model.LOGIN_FAILED,
		})
		return nil
	}

	if !password.Validate(pas, user.Password) {
		result.Failure(&errors.SimpleBizError{
			Code: model.AUTH_PASSWORD_INVALID,
		})
		log.Logger.Error("password not match", zap.Any("result", result))
		return nil
	}
	if len(user.RoleIds) > 0 {
		role := &model.Role{}
		roleId := user.RoleIds[0]
		_, err = service.Db.Where(" dtd = false and id = ?", roleId).Get(role)
		if err != nil {
			log.Logger.Error("failure to get role message by role id", zap.Error(err))
			return err
		}
		user.DefaultPath = role.Path
	} else {
		user.DefaultPath = ""
	}
	//err = service.getRoles(user)
	//if err != nil {
	//	return err
	//}

	user.Password = ""
	result.Ok = true
	return nil
}

//Deprecated, format in client
func (service *UserService) getRoles(dto *model.UserDTO) (err error) {
	if len(dto.RoleIds) == 0 {
		return be.New("user has no role")
	}

	var roles []model.Role
	err = service.Db.Table("op_role").In("id", dto.Roles).Select("permissions, code").Find(&roles)
	if err != nil {
		return err
	}

	if len(roles) > 0 {
		var permissionMap = map[string]bool{}
		var permissions []string
		for _, role := range roles {
			if len(role.Permissions) > 0 {
				for _, permission := range role.Permissions {
					if !permissionMap[role.Id] {
						permissionMap[role.Id] = true
						permissions = append(permissions, permission)
					}
				}
			}
		}

		permissionArray := communal.StringArray(permissions)
		dto.Permissions = &(permissionArray)
		dto.Roles = &roles
	}

	return nil
}

func (service *UserService) GetMe(ctx context.Context, id *int64, result *communal.Result) (err error) {
	if *id <= 0 {
		result.Failure(errors.InvalidParams())
		return nil
	}

	user := result.Data.(*model.UserDTO)
	_, err = service.Db.ID(*id).Get(user)
	if err != nil {
		return err
	}

	if user.Id <= 0 {
		result.Failure(errors.NotFound())
		return nil
	}
	if len(user.RoleIds) > 0 {
		user.Roles = &[]model.Role{}
		err = service.Db.In("id", user.RoleIds).Cols("permissions", "code").Find(user.Roles)
		if err != nil {
			return err
		}
	}
	user.Password = ""
	result.Success(user)
	//log.Logger.Debug("get me by id ", zap.Int64("id", *id), zap.Any("user", user))
	return
}

func (service *UserService) Get(ctx context.Context, id *int64, result *communal.Result) (err error) {
	if *id <= 0 {
		result.Failure(errors.InvalidParams())
		return nil
	}

	user := result.Data.(*model.UserDTO)
	_, err = service.Db.ID(*id).Get(user)
	if err != nil {
		return err
	}

	if user.Id <= 0 {
		result.Failure(errors.NotFound())
		return nil
	}

	if len(user.RoleIds) > 0 {
		user.Roles = &[]model.Role{}
		err = service.Db.In("id", user.RoleIds).Find(user.Roles)
		if err != nil {
			return err
		}
	}

	user.Password = ""
	result.Success(user)
	//log.Logger.Debug("get user by id ", zap.Int64("id", *id), zap.Any("user", user))
	return
}

func (service *UserService) Filter(ctx context.Context, filter *model.UserFilter, result *communal.FilterResult) (err error) {
	var users []model.UserDTO
	ss := service.Db.Where("dtd = false")
	if filter.Keyword != "" {
		ss.And("username like ? or mobile = ?", "%"+filter.Keyword+"%", filter.Keyword)
	}
	if filter.OrgId > 0 {
		ss.And("org_id=?", filter.OrgId)
	}

	if len(filter.Roles) > 0 {
		roles := communal.StringArray(filter.Roles)
		bts, err := (&roles).ToDB()
		if err != nil {
			return err
		}
		ss.And("role_ids @> ?", string(bts))
	}
	ss.Limit(filter.Limit(), filter.Skip())
	var count int64
	if filter.Scope == "min" {
		ss.Cols("id", "nickname")
	} else {
		ss.Omit("password")
	}
	count, err = ss.FindAndCount(&users)

	if err != nil {
		return err
	}
	var allRoles []model.Role
	err = service.Db.Where("dtd = false").Find(&allRoles)
	rolesMap := make(map[string]model.Role)
	for _, role := range allRoles {
		rolesMap[role.Id] = role
	}
	if filter.Scope != "min" {
		//todo select by single SQL
		for i := 0; i < len(users); i++ {
			userRoles := make([]interface{}, 0)
			for j := 0; j < len(users[i].RoleIds); j++ {
				//users[i].Roles = &[]model.Role{}
				//err = service.Db.In("id", users[i].RoleIds).Find(users[i].Roles)
				//if err != nil {
				//	return err
				//}
				userRoles = append(userRoles, rolesMap[users[i].RoleIds[j]])
			}
			roles := []model.Role{}
			rolesBytes, err := json.Marshal(userRoles)
			err = json.Unmarshal(rolesBytes, &roles)
			if err != nil {
				continue
			}
			users[i].Roles = &roles
		}
	}

	// 机构
	if filter.ExistOrg {
		var oids []int64
		for i, _ := range users {
			oids = append(oids, users[i].OrgId)
		}

		var orgs []model.JnOrg
		err = service.Db.Table("org").Where("dtd = false").In("id", oids).Find(&orgs)
		if err != nil {
			return err
		}
		orgMap := make(map[int64]model.JnOrg)
		for _, org := range orgs {
			orgMap[org.Id] = org
		}

		for i, _ := range users {
			if orgMap[users[i].OrgId].Id > 0 {
				users[i].OrgName = orgMap[users[i].OrgId].Name
			}
		}

	}
	result.Data = &users
	log.Logger.Debug("get user list", zap.Any("user", users))
	result.Page = filter.GetPager(count)
	result.Ok = true
	return
}

func (service *UserService) Insert(ctx context.Context, form *model.UserDTO, result *communal.Result) error {
	if form.Password == "" {
		result.Failure(&errors.SimpleBizError{
			Code: model.AUTH_PASSWORD_INVALID,
		})
		log.Logger.Error("password not match", zap.Any("result", result))
		return nil
	}
	if form.Mcc == "" {
		form.Mcc = model.MCC_DEFAULT
	}

	//去掉前后空格
	form.Username = strings.TrimSpace(form.Username)
	pas := strings.TrimSpace(form.Password)
	form.Password, _ = password.Encrypt(pas)

	form.InitBaseFields()

	_, err := service.Db.Insert(form)
	if err != nil {
		return err
	}

	result.Ok = true
	return nil
}

func (service *UserService) Update(ctx context.Context, form *model.UserDTO, result *communal.Result) error {
	if form.Id <= 0 {
		result.Failure(errors.InvalidParams().AddError(errors.InvalidField("id", errors.FIELD_BAD_FORMAT, "")))
		return nil
	}
	form.Username = strings.TrimSpace(form.Username)
	if len(form.Path) > 0 {
		form.OrgId = form.Path[len(form.Path)-1]
	}
	form.Lut = time.Now()
	_, err := service.Db.ID(form.Id).Omit("crt", "password").Update(form)
	if err != nil {
		return err
	}

	result.Ok = true
	return nil
}

func (service *UserService) Delete(ctx context.Context, id int64, result *communal.Result) error {
	if id < 1 {
		result.Failure(errors.NotFound())
		return nil
	}
	user := &model.User{}
	user.Dtd = true
	user.Lut = time.Now()

	if _, err := service.Db.ID(id).Cols("dtd", "lut").Update(user); err != nil {
		return err
	}
	result.Ok = true
	return nil
}

func (service *UserService) ResetPassword(ctx context.Context, form *model.User, result *communal.Result) error {
	if form.Id < 1 || form.Password == "" {
		result.Failure(errors.NotFound())
		return nil
	}
	form.Password, _ = password.Encrypt(form.Password)
	form.Lut = time.Now()
	_, err := service.Db.ID(form.Id).Cols("password", "lut").Update(form)
	if err != nil {
		return err
	}
	result.Ok = true
	return nil
}

func (service *UserService) UpdatePassword(ctx context.Context, form *model.RestPassword, result *communal.Result) error {
	log.Logger.Debug("update password")
	user := result.Data.(*model.User)
	userId := form.Id

	_, err := service.Db.ID(userId).Get(user)
	if err != nil {
		log.Logger.Error("failed to get user by mobile", zap.Any("err", err))
		return err
	}
	if user.Id <= 0 {
		return err
	}
	if !password.Validate(form.OldPassword, user.Password) {
		result.Failure(&errors.SimpleBizError{
			Code: model.AUTH_PASSWORD_INVALID,
		})
		log.Logger.Error("password not match", zap.Any("result", result))
		return err
	}
	form.NewPassword, _ = password.Encrypt(form.NewPassword)
	user.Password = form.NewPassword
	_, err = service.Db.ID(userId).Cols("password").Update(user)
	if err != nil {
		return err
	}
	result.Ok = true
	return nil
}
