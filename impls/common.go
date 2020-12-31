package impls

import (
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/astaxie/beego/httplib"
)

type Common struct {
	UserNameRegx *regexp.Regexp
	PasswordRegx *regexp.Regexp
	EmailRegx    *regexp.Regexp
	MobileRegx   *regexp.Regexp
}

var CommonSrv Common

func init() {
	CommonSrv.Init()
}

func (this *Common) Init() {
	//用户名正则，4到16位（字母，数字，下划线，减号）
	this.UserNameRegx = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,16}$`)
	//密码(以字母开头，长度在6~18之间，只能包含字母、数字和下划线)
	this.PasswordRegx = regexp.MustCompile(`^[a-zA-Z]\w{5,17}$`)
	//
	this.EmailRegx = regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)

	this.MobileRegx = regexp.MustCompile(`^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`)

}

func (this *Common) PwdHash(pwd string) string {
	h := md5.New()
	h.Write([]byte(pwd + "eyeits"))
	sumBytes := h.Sum(nil)
	return hex.EncodeToString(sumBytes)
}

func (this *Common) HttpReq(req *httplib.BeegoHTTPRequest) (string, int, error) {
	resp, err := req.Response()
	if err != nil {
		return "", -1, err
	}
	if resp.Body == nil {
		return "", resp.StatusCode, nil
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", resp.StatusCode, err
		}
		out, err := ioutil.ReadAll(reader)
		return string(out), resp.StatusCode, err
	}
	out, err := ioutil.ReadAll(resp.Body)
	return string(out), resp.StatusCode, err
}

func (this *Common) UrlRequest(url string, method string, username string, password string, params map[string]string, body string, headers map[string]string, connTimeout int, execTimeout int, retry int, delay int) (res string, code int, err error) {
	var req *httplib.BeegoHTTPRequest
	if len(method) == 0 {
		method = "get"
	}

	switch method {
	case "get":
		req = httplib.Get(url)
	case "post":
		req = httplib.Post(url)
	case "delete":
		req = httplib.Delete(url)
	case "put":
		req = httplib.Put(url)
	case "head":
		req = httplib.Head(url)
	default:
		return "", -1, errors.New("Invalid http method:" + method)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header(k, v)
		}

	}

	if params != nil {
		for k, v := range params {
			req.Param(k, v)
		}
	}

	if len(body) > 0 {
		req.Body(body)
	}

	if len(username) > 0 {
		req.SetBasicAuth(username, password)
	}

	defaultSetting := httplib.BeegoHTTPSettings{
		UserAgent:        "eyeits",
		ConnectTimeout:   time.Duration(connTimeout) * time.Second,
		ReadWriteTimeout: time.Duration(execTimeout) * time.Second,
		Gzip:             true,
		DumpBody:         true,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	req.Setting(defaultSetting)

	if retry < 0 {
		retry = 0
	}
	retry++
	for i := 0; i < retry; i++ {
		res, code, err = this.HttpReq(req)
		if err == nil {
			break
		}
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	return
}

func (this *Common) InternalRequest(api string, body string) (res string, code int, err error) {
	masterApiAddr, err := GlobalConfig.GetMasterApiAddr()
	if err != nil {
		return "", -1, err
	}
	return this.UrlRequest(masterApiAddr+api, "post", "", "", nil, body, nil, 15, 30, 1, 1)
}
