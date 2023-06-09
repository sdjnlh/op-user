package service

import (
	"context"
	"github.com/sdjnlh/communal/errors"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/model/module"
	"time"

	"github.com/sdjnlh/communal"
)

type GroupService struct {
	*communal.Module
}

var Group = &GroupService{
	Module: module.Group,
}

func (service *GroupService) Get(ctx context.Context, id *int64, receiver *communal.Result) (err error) {
	receiver.Data = &model.Group{}
	return service.Module.Get(ctx, id, receiver)
}

func (service *GroupService) List(ctx context.Context, filter *model.GroupFilter, receiver *communal.FilterResult) (err error) {
	receiver.Data = &[]model.Group{}
	return service.Db.Where("dtd=false").Find(receiver.Data)
}

func (service *GroupService) Create(ctx context.Context, form *model.Group, result *communal.Result) (err error) {
	return service.Module.Create(ctx, form, result)
}

func (service *GroupService) Update(ctx context.Context, form *model.Group, result *communal.Result) (err error) {
	_, err = service.Db.ID(form.Id).Update(form)
	return
}

func (service *GroupService) Delete(ctx context.Context, id *int64, result *communal.Result) error {
	if *id < 1 {
		result.Failure(errors.NotFound())
		return nil
	}
	group := &model.Group{
		Dtd: true,
		Lut: time.Now(),
	}
	if _, err := service.Db.ID(id).Cols("dtd", "lut").Update(group); err != nil {
		return err
	}
	result.Ok = true
	return nil
}
