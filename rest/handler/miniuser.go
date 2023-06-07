package handler

import (
	"code.letsit.cn/go/common"
	"code.letsit.cn/go/common/id"
	"code.letsit.cn/go/common/web"
	"code.letsit.cn/go/op-user/model"
	"code.letsit.cn/go/op-user/service"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/muesli/cache2go"
	"io"
	"os"
	"strconv"
	"time"
)

var MyCoche *cache2go.CacheTable

type myStruct struct {
	text     string
	moreData []byte
}

type MiniUserApi struct {
	*web.RestHandler
}

func NewMiniUserApi() *MiniUserApi {
	return &MiniUserApi{
		RestHandler: web.DefaultRestHandler,
	}
}
func (api *MiniUserApi) get(c *gin.Context) {
	// 获取id
	form := &model.MiniUser{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {

		c.String(400, "id 参数错误")
		c.Abort()
		return
	}
	form.Id = id
	// 调用服务方法获取数据库记录
	if err := service.MiniUser.Get(form); err != nil {
		c.String(500, "获取记录失败")
		c.Abort()
		return
	}
	c.JSON(200, form)
}

func (api *MiniUserApi) update(c *gin.Context) {
	form := &model.MiniUser{}
	if err := c.Bind(form); err != nil {
		c.String(400, "表单参数错误")
		c.Abort()
		return
	}
	// 调用更新方法
	if err := service.MiniUser.Update(form); err != nil {
		c.String(500, "更新失败")
		c.Abort()
		return
	}
	c.JSON(200, form)
}

// get 获取单个记录
func (api *MiniUserApi) list(c *gin.Context) {
	// 绑定表单属性和分页属性

	form := &model.MiniUser{}
	page := &common.Page{}

	list := &[]model.MiniUser{}
	if err := c.Bind(page); err != nil {
		c.String(400, "更新失败")
		c.Abort()
		return
	}
	if err := service.MiniUser.List(form, page, list); err != nil {
		c.String(500, "获取记录失败")
		c.Abort()
		return
	}

	// 返回结果
	r := make(map[string]interface{})
	r["data"] = list
	r["page"] = page
	c.JSON(200, r)
}

// delete 删除单个记录
func (api *MiniUserApi) delete(c *gin.Context) {
	// 获取id
	form := &model.MiniUser{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(400, "id 参数错误")
		c.Abort()
		return
	}
	form.Id = id
	//删除单个记录
	if err := service.MiniUser.Delete(form); err != nil {
		c.String(500, "删除失败")
		c.Abort()
		return
	}
	c.String(200, "删除成功")
}

// save 增加单个记录
func (api *MiniUserApi) save(c *gin.Context) {
	form := &model.MiniUser{}
	// 绑定表单数据
	if err := c.Bind(form); err != nil {
		c.String(400, "id 参数错误")
		c.Abort()
		return
	}
	// 保存成为新的记录
	if err := service.MiniUser.Save(form); err != nil {
		c.String(500, "保存失败")
		c.Abort()
		return
	}
	// 保存成功后返回新的记录
	c.JSON(200, form)
}
func (api *MiniUserApi) bindUserList(c *gin.Context) {
	miniUser := &model.MiniUser{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(400, "id 参数错误")
		c.Abort()
		return
	}
	miniUser.BdId = id
	r := &map[string]interface{}{}
	if err := service.MiniUser.GetBindUserList(miniUser, r); err != nil {
		c.String(500, "获取失败")
		c.Abort()
		return
	}
	c.JSON(200, r)

}

func (api *MiniUserApi) login(c *gin.Context) {
	form := &model.MiniUser{}
	// 绑定表单数据
	if err := c.Bind(form); err != nil {
		c.String(400, "id 参数错误")
		c.Abort()
		return
	}
	// 保存成为新的记录
	if err := service.MiniUser.Login(form); err != nil {
		c.String(500, "保存失败")
		c.Abort()
		return
	}
	// 保存成功后返回新的记录
	c.JSON(200, form)
}
func (api *MiniUserApi) decodePhone(c *gin.Context) {
	decode := &model.DecodePhone{}
	if err := c.Bind(decode); err != nil {
		c.String(400, "参数错误")
		c.Abort()
	}
	sessionKey, err := base64.StdEncoding.DecodeString(decode.SessionKey)
	iv, err := base64.StdEncoding.DecodeString(decode.Iv)
	decodeBytes, err := base64.StdEncoding.DecodeString(decode.EncryptedData)
	if err != nil {
		c.String(500, "内部服务器错误")
		c.Abort()
		return
	}
	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		c.String(500, "内部服务器错误")
		c.Abort()
		return
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeBytes))
	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)

	phone := &model.Phone{}
	if err := json.Unmarshal(origData, phone); err != nil {
		c.String(500, "内部服务器错误")
		c.Abort()
		return
	}
	//将解析获取到的数据跟新入数据库，并返回信息
	if _, err := service.MiniUser.UpdatePhone(decode.OpenId, phone.PhoneNumber); err != nil {
		c.String(500, "内部服务器错误")
		c.Abort()
		return
	} else {
		c.JSON(200, phone)
	}

}

//小程序获取校验码
func (api *MiniUserApi) SendPhoneCode(c *gin.Context) {
	form := &model.PhoneCode{}
	if err := c.Bind(form); err != nil {
		c.String(400, "参数错误")
		c.Abort()
		return
	}
	if err, bool := service.User.CheckMobileUser(form); err != nil {
		c.String(400, "checkMobile失败")
		c.Abort()
		return
	} else {
		if bool {
			//发送短信验证码
			code, has := service.Sender.SendCode(form.PhoneNumber)
			if !has {
				c.String(400, "获取手机号失败")
				c.Abort()
				return
			}
			val := myStruct{code, []byte{}}
			MyCoche.Add(form.PhoneNumber, 600*time.Second, &val)

		} else {
			r := map[string]interface{}{}
			r["data"] = "获取校验码失败!"
			c.JSON(200, r)
			c.Abort()
			return
		}

	}
}

// 小程序手机号校验登录
func (api *MiniUserApi) WxCodeVerifyUser(c *gin.Context) {
	form := &model.PhoneCode{}
	if err := c.Bind(form); err != nil {
		c.String(400, "参数错误")
		c.Abort()
		return
	}
	r := map[string]interface{}{}
	val, err := MyCoche.Value(form.PhoneNumber)
	if err != nil {
		c.String(400, "服务器错误")
		c.Abort()
		return
	}
	mystruct := val.Data().(*myStruct)
	if mystruct.text == form.Code {
		//获取user信息
		user := &model.MiniUser{}
		if err := service.User.GetUser(form.PhoneNumber, user); err != nil {
			c.String(500, "获取加盟商失败！")
			c.Abort()
			return
		}
		if user.Id != 0 {
			r["data"] = user
			r["status"] = true
			c.JSON(200, r)
		} else {
			r["status"] = false
			c.JSON(200, r)
		}

	}
}

func (api *MiniUserApi) uploadImage(c *gin.Context) {
	defer c.Request.Body.Close()
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(400, "图片id 参数错误")
		c.Abort()
		return
	}

	ids, _ := id.Next()
	userImageId := strconv.FormatInt(ids, 10)
	userId, _ := c.GetPostForm("userId")

	pathstr := "/topiot/fs/up"
	path := ".." + pathstr
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			c.String(400, "创建文件失败")
			c.Abort()
			return
		} //0777也可以os.ModePerm
		if err := os.Chmod(path, os.ModePerm); err != nil {
			c.String(400, "chmod失败")
			c.Abort()
			return
		}
	}

	urlpath := "/" + userImageId + ".png"
	str := path + "/" + userImageId + ".png"
	fWrite, err := os.Create(str)
	//将此路径写入数据库
	if userId != "" {
		if err := service.User.InsertImagePath(userId, urlpath); err != nil {
			c.String(450, "chmod失败")
			c.Abort()
			return
		}
	}
	if _, err := io.Copy(fWrite, file); err != nil {
		//logger.Info("文件保存失败")
	}

	defer fWrite.Close()

	r := make(map[string]interface{})
	r["data"] = urlpath
	c.JSON(200, r)
	c.Abort()

	return
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func init() {
	MyCoche = cache2go.Cache("myCache")
}
func (api *MiniUserApi) Register(r gin.IRouter) {
	r.POST("/v1/mini-user/login", api.login) //c端登陆
	r.GET("/v1/mini-user/:id", api.get)
	r.DELETE("/v1/mini-user/:id", api.delete) // 删除单条记录
	r.PUT("/v1/mini-user/:id", api.update)
	r.GET("/v1/mini-users", api.list) // 分页查询多条记录
	r.POST("/v1/mini-user", api.save) // 增加单条记录
	// 获取r该用户下邀请绑定的用户
	r.GET("/v1/mini-user/:id/bind", api.bindUserList)
	r.POST("/v1/mini-user/phone_number", api.decodePhone) //获取手机号
	// 验证码
	r.POST("/v1/mini-user/wxLogin/phone", api.SendPhoneCode)
	r.POST("/v1/mini-user/wxSender/phoneVerify", api.WxCodeVerifyUser)
	//图片
	r.POST("/v1/mini-user/uploadimage", api.uploadImage)

}
