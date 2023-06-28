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

type OrgService struct {
	*communal.Module
}

var Org = &OrgService{
	Module: module.Org,
}

func (service *OrgService) Get(ctx context.Context, id *string, result *communal.Result) (err error) {
	org := result.Data.(*model.Org)
	_, err = service.Db.Table("organization").ID(*id).Get(org)
	if org.Id == 0 || err != nil {
		result.Failure(errors.NotFound())
		return nil
	}

	result.Ok = true
	return
}

func (service *OrgService) Filter(ctx context.Context, filter *model.OrgFilter, result *communal.FilterResult) (err error) {
	orgs := result.Data.(*[]model.Org)
	ss := service.Db.Where("dtd = false")
	log.Logger.Debug("keyword text", zap.Any("Keywork text", filter.Name))
	if filter.Name != "" {
		ss.And("name like ?", "%"+filter.Name+"%")
	}
	count, err := ss.Table("organization").Limit(filter.Limit(), filter.Skip()).FindAndCount(orgs)
	if err != nil {
		return err
	}
	result.Page = filter.GetPager(count)
	result.Ok = true
	return
}

func (service *OrgService) Insert(ctx context.Context, org *model.Org, result *communal.Result) error {
	//org.Id = org.Name
	org.InitBaseFields()
	_, err := service.Db.Table("organization").Insert(org)
	if err != nil {
		return err
	}
	result.Ok = true
	return nil
}

func (service *OrgService) Update(ctx context.Context, org *model.Org, result *communal.Result) (err error) {
	org.Lut = time.Now()
	_, err = service.Db.Table("organization").Where("id=?", org.Id).Omit("crt").Update(org)
	if err != nil {
		return err
	}

	result.Ok = true
	return nil
}
func (service *OrgService) Delete(ctx context.Context, id string, result *communal.Result) error {
	org := &model.Org{}
	org.Dtd = true
	org.Lut = time.Now()
	_, err := service.Db.Table("organization").Where("id=?", id).Cols("dtd", "lut").Update(org)
	if err != nil {
		return err
	}
	result.Ok = true
	return nil
}
