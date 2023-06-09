package handler

import (
	"fmt"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/db"
	"github.com/sdjnlh/communal/errors"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/rpc"
	"github.com/sdjnlh/communal/web"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/opu"
	"github.com/sdjnlh/op-user/service"
	senderModel "github.com/sdjnlh/sender/model"
	senderService "github.com/sdjnlh/sender/service"
	"math/rand"
	"time"

	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type UserAPI struct {
	*web.RestHandler
}

func NewUserAPI() *UserAPI {
	return &UserAPI{
		RestHandler: web.DefaultRestHandler,
	}
}

func (api *UserAPI) Logout(c *gin.Context) {
	//session.ClearUserSession(c)
}

func (api *UserAPI) Me(c *gin.Context) {
	log.Logger.Debug("Me")
	uid := api.UID(c)

	dto := model.UserDTO{}
	result := &communal.Result{
		Data: &dto,
	}

	err := service.User.GetMe(context.Background(), &uid, result)

	//log.Logger.Debug("get  user", zap.Any("me", dto))
	result.Ok = true
	api.ResultWithError(c, result, err)
}

func (api *UserAPI) Login(c *gin.Context) {
	var form = model.User{}
	err := api.Bind(c, &form)
	if err != nil {
		log.Logger.Debug("fail to bind login params", zap.Error(err))
		api.BadRequestWithError(c, err)
		return
	}

	user := model.UserDTO{}
	result := &communal.Result{
		Data: &user,
	}

	err = service.User.Login(context.Background(), &form, result)
	if err != nil || !result.Ok {
		api.ResultWithError(c, result, err)
		return
	}

	//log.Logger.Debug("user ", zap.Any(":", user))
	err = opu.Api.StateManager.SetUser(c, user.State())
	if err != nil {
		api.ResultWithError(c, result, err)
		return
	}
	c.SetCookie(communal.UserIdKey, strconv.FormatInt(user.Id, 10), opu.Api.StateManager.MaxAge(), opu.Api.StateManager.Path(), opu.Api.StateManager.Domain(), false, false)
	c.SetCookie(communal.UserNicknameKey, user.Nickname, opu.Api.StateManager.MaxAge(), opu.Api.StateManager.Path(), opu.Api.StateManager.Domain(), false, false)
	c.SetCookie(communal.UserRoleKey, strings.Join(user.RoleIds, ","), opu.Api.StateManager.MaxAge(), opu.Api.StateManager.Path(), opu.Api.StateManager.Domain(), false, false)
	c.SetCookie(communal.UserOrgIdKey, strconv.FormatInt(user.OrgId, 10), opu.Api.StateManager.MaxAge(), opu.Api.StateManager.Path(), opu.Api.StateManager.Domain(), false, false)
	api.SuccessWithData(c, result)
}

func (api *UserAPI) Filter(c *gin.Context) {
	var form = &model.UserFilter{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	//strOrg, _ := c.Cookie(communal.UserOrgIdKey)
	//orgid, _ := strconv.ParseInt(strOrg, 10, 64)
	//if form.OwnerId <= 0 {
	//	form.OwnerId = orgid
	//}
	//if form.OrgId <= 0 {
	//	form.OrgId = orgid
	//}

	result := communal.NewFilterResult(&[]model.UserDTO{})
	err = service.User.Filter(context.Background(), form, result)

	//log.Logger.Debug("list user", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

//更新密码
func (api *UserAPI) UpdatePassword(c *gin.Context) {
	var form = &model.RestPassword{}
	err := api.Bind(c, form)
	form.Id = api.UID(c)
	log.Logger.Debug("update password form :", zap.Any("update password form ", form))
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	result := &communal.Result{
		Data: &model.User{},
	}
	err = service.User.UpdatePassword(context.Background(), form, result)
	log.Logger.Debug("list user", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *UserAPI) ResetPassword(c *gin.Context) {
	var form = &model.User{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	result := &communal.Result{
		Data: &model.User{},
	}

	err = service.User.ResetPassword(context.Background(), form, result)

	log.Logger.Debug("reset password", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *UserAPI) Update(c *gin.Context) {
	var form = &model.UserDTO{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	strOrg, _ := c.Cookie(communal.UserOrgIdKey)
	orgid, _ := strconv.ParseInt(strOrg, 10, 64)
	form.OrgId = orgid
	result := &communal.Result{}
	err = service.User.Update(context.Background(), form, result)

	log.Logger.Debug("update user", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *UserAPI) Delete(c *gin.Context) {
	id, err := api.ValidateInt64Id(c)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	result := &communal.Result{}

	err = service.User.Delete(context.Background(), id, result)

	log.Logger.Debug("delete user", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *UserAPI) Get(c *gin.Context) {
	id, err := api.ValidateInt64Id(c)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	result := &communal.Result{Data: &model.UserDTO{}}
	err = service.User.Get(context.Background(), &id, result)

	//log.Logger.Debug("get user", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

func (api *UserAPI) Insert(c *gin.Context) {
	var form = &model.UserDTO{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	//strOrg, _ := c.Cookie(communal.UserOrgIdKey)
	//orgid, _ := strconv.ParseInt(strOrg, 10, 64)
	//form.OrgId = orgid
	result := &communal.Result{}
	if opu.Api.Rpc {
		err = rpc.Call(context.Background(), "OpUser", "Insert", form, result)
		log.Logger.Error("failed to call rpc", zap.Any("error", err))
	} else {
		err = service.User.Insert(context.Background(), form, result)
	}

	//log.Logger.Debug("create user", zap.Any("result", result))
	api.ResultWithError(c, result, err)
}

//发送code到邮件
func (api *UserAPI) SenderEmail(c *gin.Context) {
	var form = &model.ForgetForm{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	has := db.SetKeyWithExpireTime(opu.Service.Redis, form.Email, code, model.CaptchaTimeout)
	if !has {
		log.Logger.Info("redis err")
	}
	var msg = &senderModel.Message{
		Type:       "email",
		Recipient:  form.Email,
		TemplateId: "verify-code",
		Params:     make(map[string]interface{}),
	}
	msg.Params["subject"] = "Thanks for your registration"
	msg.Params["code"] = code
	if err := senderService.SendDenialLetter(msg); err != nil {
		log.Logger.Error("failed to send email", zap.Any("error", err))
		return
	}
	api.Success(c)
}

//验证code信心
func (api *UserAPI) VerifyCode(c *gin.Context) {
	var form = &model.ForgetForm{}
	result := &communal.Result{}
	err := api.Bind(c, form)
	if err != nil {
		api.BadRequestWithError(c, err)
		return
	}
	if form.Code != "" {
		code, errRedis := db.RedisGet(opu.Service.Redis, form.Email)
		if errRedis != nil {
			log.Logger.Info("error redis")
		}
		if form.Code == code {
			result.Ok = true
			log.Logger.Info("验证通过！")
		} else {
			err := &errors.SimpleBizError{
				Code: model.LOGIN_IDENTITY_INVALID,
			}
			result.Failure(err)

			log.Logger.Info("验证失败！")
		}
	}
	api.Result(c, result)

}
func (api *UserAPI) Register(router gin.IRouter) {
	v1 := router.Group("/v1")
	v1.POST("/login", api.Login)
	v1.GET("/me", web.UserInterceptor, api.Me)
	v1.POST("/logout", api.Logout)
	v1.GET("/users", web.UserInterceptor, api.Filter)
	v1.GET("/users/:id", web.UserInterceptor, api.Get)
	v1.POST("/users", web.UserInterceptor, api.Insert)
	v1.PUT("/password/reset", web.UserInterceptor, api.ResetPassword)
	v1.PUT("/password", web.UserInterceptor, api.UpdatePassword)
	v1.PUT("/users", web.UserInterceptor, api.Update)
	v1.PUT("/users/:id", web.UserInterceptor, api.Update)
	v1.DELETE("/users/:id", web.UserInterceptor, api.Delete)
	v1.POST("/users/senderemail", api.SenderEmail) // 邮箱验证码
	v1.POST("/users/getverifycode", api.VerifyCode)
}
