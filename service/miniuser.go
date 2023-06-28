package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/model/module"

	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/log"
	"go.uber.org/zap"
)

// 服务模板
type miniUserService struct {
	*communal.Module
}

var MiniUser = &miniUserService{
	Module: module.MiniUser,
}

func (service *miniUserService) UpdatePhone(openId string, phoneNumber string) (bool, error) {
	miniUser := &model.MiniUser{}
	//先查询是否有该手机号的的用户
	wxUser := &model.MiniUser{}
	if _, err := service.Db.Where("mobile=? and dtd <> true", phoneNumber).Get(miniUser); err != nil {
		return false, err
	}
	if _, err := service.Db.Where("openid = ? and dtd <> true", openId).Get(wxUser); err != nil {
		return false, err
	}
	if wxUser.Mobile != "" && wxUser.Password != "" {
		return false, nil
	}
	if miniUser.Id != 0 {
		wxUser.Password = miniUser.Password
		if miniUser.Region != "" && wxUser.Region == "" {
			wxUser.Region = miniUser.Region
		}
		if miniUser.Company != "" && wxUser.Company == "" {
			wxUser.Company = miniUser.Company
		}
		if miniUser.Username != "" && wxUser.Username == "" {
			wxUser.Username = miniUser.Username
		}
		miniUser.Dtd = true
		if _, err := service.Db.ID(miniUser.Id).MustCols("dtd").Update(miniUser); err != nil {
			return false, err
		}
	}
	wxUser.Mobile = phoneNumber
	if _, err := service.Db.ID(wxUser.Id).MustCols("mobile", "password").Update(wxUser); err != nil {
		return false, err
	}

	//查询数据将数据返回到前台，所需要的数据
	return false, nil
}
func (service *miniUserService) Login(form *model.MiniUser) error {

	wxlogin := model.Wxlogin
	wxappId := model.AppId
	appsecret := model.Secret
	resp, err := http.Get(fmt.Sprintf(wxlogin, wxappId, appsecret, form.Code))
	if err != nil {

		return err
	}
	response, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return err
	}

	if err != nil {
		return err
	}
	if err = json.Unmarshal(response, form); err != nil {
		return err
	}
	//获取accessToken
	token, err := service.RequestToken(wxappId, appsecret)
	if err != nil {
		return err
	}
	form.Token = token

	// 根据token 查询用户

	newUser := &model.MiniUser{}

	if _, err := service.Db.Where("openid = ? and dtd<>true ", form.Openid).Get(newUser); err != nil {
		return err
	}
	if newUser.Id == 0 {
		// 用户不存在  新建用户再返回
		if err := service.Save(form); err != nil {
			return err
		}
	} else {
		form.Id = newUser.Id
		form.Mobile = newUser.Mobile
		form.AvatarUrl = newUser.AvatarUrl
	}
	return nil
}

func (service *miniUserService) RequestToken(appid, secret string) (string, error) {
	u, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		return "", err
	}
	paras := &url.Values{}
	//设置请求参数
	paras.Set("appid", appid)
	paras.Set("secret", secret)
	paras.Set("grant_type", "client_credential")
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return "", err
	}
	if jMap["errcode"] == nil || jMap["errcode"] == 0 {
		accessToken, _ := jMap["access_token"].(string)
		return accessToken, nil
	} else {
		return "", errors.New(jMap["code"].(string) + ":" + jMap["msg"].(string))
	}
}

// Get 获取单个用户
func (service *miniUserService) Get(form *model.MiniUser) error {
	// 更新数据库中的记录
	if _, err := service.Db.ID(form.Id).Get(form); err != nil {
		return err
	}
	return nil
}

// list 获取多个项目列表
func (service *miniUserService) List(form *model.MiniUser, page *communal.Page, list *[]model.MiniUser) error {
	// 分页查询
	sql := service.Db.NewSession()
	defer sql.Commit()
	//if page.K != "" {
	//	sql.Where("name like ?", "%"+page.K+"%")
	//}
	if cnt, err := sql.Limit(page.Limit(), page.Skip()).Desc("crt").FindAndCount(list, form); err != nil {
		fmt.Println(err)
		return err
	} else {
		page = page.GetPager(cnt)
	}

	return nil
}

// 更新单个用户
func (service *miniUserService) Update(form *model.MiniUser) error {

	if _, err := service.Db.ID(form.Id).Update(form); err != nil {
		return err
	}

	return nil
}

// Delete 删除记录
func (service *miniUserService) Delete(form *model.MiniUser) error {
	// 删除记录
	if form.Id < 1 {
		return nil
	}
	form.Dtd = true
	form.Lut = time.Now()
	if _, err := service.Db.ID(form.Id).Cols("dtd", "lut").Update(form); err != nil {
		return err
	}

	return nil
}

// Save 保存记录
func (service *miniUserService) Save(form *model.MiniUser) error {
	form.InitBaseFields()
	if _, err := service.Db.Insert(form); err != nil {
		log.Logger.Error("fail to save item", zap.Any("err", err))
		return err
	}

	return nil
}
func (service *miniUserService) GetBindUserList(form *model.MiniUser, r *map[string]interface{}) error {
	miniUsers := &[]model.MiniUser{}
	if err := service.Db.Where("bd_id=?", form.BdId).Find(miniUsers); err != nil {
		return err
	}
	(*r)["data"] = miniUsers
	return nil
}

// check user
func (service *UserService) CheckMobileUser(phoneCode *model.PhoneCode) (error, bool) {
	newUser := &model.MiniUser{}
	if _, err := service.Db.Table("mini_user").Where("mobile=? or id=? and dtd<>true", phoneCode.PhoneNumber, phoneCode.UserId).Get(newUser); err != nil {
		return err, false
	}
	// 没有
	if newUser.Id == 0 {
		form := &model.MiniUser{}
		form.InitBaseFields()
		form.Dtd = false
		form.Mobile = phoneCode.PhoneNumber
		if _, err := service.Db.Table("mini_user").Cols("id", "crt", "lut", "password", "mobile", "dtd", "bd_id").Insert(form); err != nil {
			log.Logger.Error("create user failed", zap.Any("err", err))
			return err, false
		}
		return nil, true
	} else {
		//有
		if phoneCode.MobileFilter != "" {
			// user用户更新手机号
			newUser.Mobile = phoneCode.PhoneNumber
			if _, err := service.Db.ID(newUser.Id).Cols("mobile").Update(newUser); err != nil {
				return err, false
			}
		}
		return nil, true

	}
	return nil, false
}

func (service *UserService) GetUser(mobile string, user *model.MiniUser) error {
	if _, err := service.Db.Table("mini_user").Where("mobile=? and dtd<>true", mobile).Get(user); err != nil {
		return err
	}
	return nil
}

func (service *UserService) InsertImagePath(userId string, path string) (err error) {
	//查询shop和mer_shop
	user := &model.MiniUser{}
	if _, err := service.Db.ID(userId).Get(user); err != nil {
		return err
	}
	// 处理图片路径的双引号
	path = strings.ReplaceAll(path, "\"", "")
	if path != "" {
		user.AvatarUrl = path
	}

	if _, err := service.Db.ID(userId).Cols("avatar_url").Update(user); err != nil {
		return err
	}

	return nil
}
