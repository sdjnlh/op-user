package service

import (
	"context"
	"time"

	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/errors"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/model/module"
	"go.uber.org/zap"
)

type RoleService struct {
	*communal.Module
}

var Role = &RoleService{
	Module: module.Role,
}

func (service *RoleService) Get(ctx context.Context, id *string, result *communal.Result) (err error) {
	role := result.Data.(*model.RoleDTO)
	_, err = service.Db.Where("id=?", *id).Get(role)
	if role.Id == "" || err != nil {
		result.Failure(errors.NotFound())
		return nil
	}

	result.Ok = true
	return
}

func (service *RoleService) Filter(ctx context.Context, filter *model.RoleFilter, result *communal.FilterResult) (err error) {
	roles := result.Data.(*[]model.Role)
	ss := service.Db.Where("dtd = false")
	log.Logger.Debug("keyword text", zap.Any("Keywork text", filter.Keyword))
	if filter.Keyword != "" {
		ss.And("name like ?", "%"+filter.Keyword+"%")
	}
	//if filter.OwnerId > 0 {
	//	ss.And("owner_id=?", filter.OwnerId)
	//}
	count, err := ss.Limit(filter.Limit(), filter.Skip()).FindAndCount(roles)
	if err != nil {
		return err
	}
	result.Page = filter.GetPager(count)
	result.Ok = true
	return
}

func (service *RoleService) Insert(ctx context.Context, role *model.RoleDTO, result *communal.Result) error {
	// role.Id = role.Name
	now := time.Now()
	role.Crt = now
	role.Lut = now
	_, err := service.Db.Insert(role)
	if err != nil {
		log.Logger.Error("create role", zap.Any("result", result))
		return err
	}
	result.Ok = true
	return nil
}

func (service *RoleService) Update(ctx context.Context, role *model.RoleDTO, result *communal.Result) (err error) {
	role.Lut = time.Now()
	_, err = service.Db.Where("id=?", role.Id).Omit("crt").Update(role)
	if err != nil {
		return err
	}

	result.Ok = true
	return nil
}
func (service *RoleService) Delete(ctx context.Context, id string, result *communal.Result) error {
	if id == "" {
		result.Failure(errors.NotFound())
		return nil
	}
	role := &model.Role{}
	role.Dtd = true
	role.Id = id
	role.Lut = time.Now()
	_, err := service.Db.Where("id=?", id).Cols("dtd", "lut").Update(role)
	if err != nil {
		return err
	}
	result.Ok = true
	return nil
}
