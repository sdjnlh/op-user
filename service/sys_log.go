package service

import (
	"code.letsit.cn/go/common"
	"code.letsit.cn/go/common/errors"
	"code.letsit.cn/go/common/log"
	"code.letsit.cn/go/op-user/model"
	"code.letsit.cn/go/op-user/model/module"
	"context"
	"go.uber.org/zap"
	"time"
)

type SysLogService struct {
	*common.Module
}

var SysLog = &SysLogService{
	Module: module.SysLog,
}

func (service *SysLogService) List(filter *model.SysLogFilter, receiver *common.FilterResult) (err error) {
	ss := service.Db.Alias("l").Where("l.dtd=false").OrderBy("crt desc")
	if filter.Username != "" {
		ss.And("l.username like '%" + filter.Username + "' ")
	}

	if !filter.StartTime.IsZero() && !filter.EndTime.IsZero() {
		s := filter.StartTime.Format("2006-01-02 03:04:05")
		e := filter.EndTime.Format("2006-01-02 03:04:05")
		ss.And("l.crt between ? and ?", s, e)
	}

	count, err := ss.Limit(filter.Limit(), filter.Skip()).FindAndCount(receiver.Data)
	if err != nil {
		return err
	}
	receiver.Page = filter.GetPager(count)
	receiver.Ok = true
	return nil
}

func (service *SysLogService) Insert(ctx context.Context, form *model.SysLog, result *common.Result) (err error) {
	if form.Id > 0 {
		if _, err = service.Db.ID(form.GetId()).UseBool("check").Update(form); err != nil {
			log.Logger.Error("fail to update item", zap.Error(err))
			return err
		}
		result.Success("")
	} else {
		form.InitBaseFields()
		err = service.Module.Create(ctx, form, result)
	}
	if err != nil {
		return err
	}
	result.Data = form
	result.Ok = true
	return err
}
func (service *SysLogService) Delete(id int64, result *common.Result) error {
	if id < 1 {
		result.Failure(errors.NotFound())
		return nil
	}
	sysLog := &model.SysLog{}
	sysLog.Dtd = true
	sysLog.Lut = time.Now()
	if _, err := service.Db.ID(id).Cols("dtd", "lut").Update(sysLog); err != nil {
		return err
	}
	result.Ok = true
	return nil
}
