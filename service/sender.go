package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var Sender sender

type sender int

var response map[string]interface{}

const (
	CpName     = "dzd"
	CpPassword = "Dms201582"
)

func (sender) SendCode(mobile string) (string, bool) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	resp, err := http.Post("http://qxt.fungo.cn/Recv_center",
		"application/x-www-form-urlencoded",
		strings.NewReader("CpName="+CpName+"&CpPassword="+CpPassword+"&DesMobile="+mobile+"&Content=【极农智能】您的验证码是"+code+",请在十分钟内完成&ExtCode=1234"))
	if err != nil {

	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	respCode, _ := response["code"]
	if "0" != respCode.(string) {
		return code, false
	}
	return code, true
}
