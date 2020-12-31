package impls

import (
	"alphawolf.com/alpha/util"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

var GfAuthSrv GfAuth

func init() {
	GfAuthSrv.Initialize()
}

type GfAuth struct {
	Url string
}

type GfLoginRes struct {
	IsLogin    string `json:"isLogin"`
	Msg        string `json:"msg"`
	LtpaToken2 string `json:"LtpaToken2"`
}

func (this *GfAuth) Initialize() {
	this.Url = beego.AppConfig.String("auth::auth_url")
}

func (this *GfAuth) Login(username string, password string) (string, error) {
	//http://tam.gf.com.cn/app/portalservice/UserServlet?method=userLoginPost&username=lihl&password=a123456&callback=
	req := httplib.Get(this.Url + "/app/portalservice/UserServlet")

	logs.Trace(this.Url + "/app/portalservice/UserServlet")

	httpSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "AlphaService",
		ConnectTimeout:   time.Duration(10) * time.Second,
		ReadWriteTimeout: time.Duration(10) * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(httpSetting)
	req.Param("method", "userLoginPost")
	req.Param("username", username)
	req.Param("password", password)

	res, _, err := CommonSrv.HttpReq(req)
	if err != nil {
		return "", err
	}
	resLen := len(res)
	if resLen > 3 {
		res = res[1 : resLen-3]
	}

	logs.Trace("res:", string(res))

	resMsgs := make([]GfLoginRes, 0)

	err = util.JsonIter.Unmarshal([]byte(res), &resMsgs)
	if err != nil {
		return "", errors.New("unmarshal login res to json failed, " + err.Error())
	}

	if len(resMsgs) < 1 {
		return "", errors.New("login res json array is null")
	}

	if resMsgs[0].IsLogin != "true" {
		return "", errors.New(resMsgs[0].Msg)
	}

	return resMsgs[0].LtpaToken2, nil
}

//<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
//<restfulResult>
//<msg>mjim¨3129.23.37!14;32;11¨3129.23.37!34;32;11</msg>
//<success>true</success>
//</restfulResult>

type DecodeTokenRes struct {
	Msg     string `xml:"msg"`
	Success bool   `xml:"success"`
}

func (this *GfAuth) DecodeToken(token string) (string, error) {
	//http://portal.gf.com.cn/app/portalservice/restful/GFCommonService/decodeLtpaToken2
	req := httplib.Post(this.Url + "/app/portalservice/restful/GFCommonService/decodeLtpaToken2")

	httpSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "AlphaService",
		ConnectTimeout:   time.Duration(10) * time.Second,
		ReadWriteTimeout: time.Duration(10) * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(httpSetting)
	req.Header("Content-Type", "application/xml")
	req.Body(token)

	res, _, err := CommonSrv.HttpReq(req)
	//logs.Info("res,code,err:",string(res),code,err.Error())
	if err != nil {
		return "", err
	}

	logs.Trace("res:", string(res))

	decodeTokenRes := &DecodeTokenRes{}
	err = xml.Unmarshal([]byte(res), decodeTokenRes)
	if err != nil {
		return "", errors.New("unmarshal decode res to xml failed, " + err.Error())
	}

	return decodeTokenRes.Msg, err
}

//public class Encryption
//{
//public String Encode(String EncodeStr)
//{
//char[] chars = EncodeStr.toCharArray();
//EncodeStr = "";
//for (int i = 0; i < chars.length; i++) {
//EncodeStr = EncodeStr + (char)(chars[i] + '\001');
//}
//return EncodeStr;
//}
//public String Decode(String DecodeStr) {
//char[] chars = DecodeStr.toCharArray();
//DecodeStr = "";
//for (int i = 0; i < chars.length; i++) {
//DecodeStr = DecodeStr + (char)(chars[i] - '\001');
//}
//return DecodeStr;
//}
//}

func (this *GfAuth) Decrypt(token string) string {
	out := make([]byte, len(token))
	for i := 0; i < len(token); i++ {
		out[i] = token[i] - 0x0001
	}

	logs.Trace("out:\n", hex.Dump(out))
	return string(out)
}

func (this *GfAuth) GetUsername(token string) string {
	members := strings.Split(token, string([]byte{0xc1, 0xa7}))
	if len(members) > 0 {
		return members[0]
	} else {
		return ""
	}
}
