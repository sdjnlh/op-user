package security

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sdjnlh/communal"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/op-user/model"
	"github.com/sdjnlh/op-user/model/module"
	"go.uber.org/zap"

	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var lch chan model.SysLog
var once sync.Once

var LogInterceptor = func(c *gin.Context) {
	v, ok := c.Get(communal.UserKey)
	if ok {
		if v.(communal.IdInf).GetId() > 0 {
			u := v.(*model.User)
			fmt.Println(u.Username)
			pa := ""
			if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
				// 上传文件的post请求
				file, header, _ := c.Request.FormFile("file")
				if file != nil {
					pa = "Upload file  " + header.Filename + strconv.FormatInt(header.Size, 10)
				} else {
					var bodyBytes []byte
					if c.Request.Body != nil {
						bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
					}
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
					// 对登录方法和创建用户的方法特殊处理 去掉敏感信息
					if c.Request.RequestURI == "/api/v1/password" || c.Request.RequestURI == "/api/v1/login" || c.Request.RequestURI == "/api/v1/users" || (strings.Index(c.Request.RequestURI, "/api/v1/users/") > -1) {
						m := make(map[string]interface{})
						if err := json.Unmarshal(bodyBytes, &m); err != nil {
							log.Logger.Error("json.Unmarshal error", zap.Any("", ""))
						}
						m["password"] = "******"
						if c.Request.RequestURI == "/api/v1/password" {
							m["OldPassword"] = "******"
							m["newPassword"] = "******"
							m["againPassword"] = "******"
						}
						bd, _ := json.Marshal(m)
						pa = string(bd)
					} else {
						pa = string(bodyBytes)
					}
				}
			} else {
				b, _ := json.Marshal(c.Params)
				pa = string(b)
			}

			l := model.SysLog{
				Ip:       getCurrentIP(c.Request),
				Uri:      c.Request.RequestURI,
				Params:   pa,
				Username: u.Nickname,
				Uid:      u.Id,
				Method:   c.Request.Method,
				Agent:    c.Request.UserAgent(),
				Ext:      nil,
			}
			l.InitBaseFields()
			lch <- l

		}
	}
	c.Next()
}

func init() {
	once.Do(func() {
		lch = make(chan model.SysLog, 999)
	})

	go func() {
		for {
			l := <-lch
			db := module.SysLog.Db
			if db == nil {
				log.Logger.Error("db == nil", zap.Any("", ""))
			}
			if _, err := db.Insert(l); err != nil {
				log.Logger.Error("db.Insert(l) error", zap.Any("", ""))
			}
		}
	}()
}

/**
		proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
*/
func getCurrentIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}
