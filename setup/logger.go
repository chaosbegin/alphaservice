package setup

import (
	"alphawolf.com/alpha/util"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func InitLogger() error {
	exePath, _ := os.Executable()
	pwd := filepath.Dir(exePath)

	exeFileName := filepath.Base(exePath)
	exeFileExt := filepath.Ext(exePath)
	if len(exeFileExt) > 0 {
		exeFileName = strings.TrimSuffix(exeFileName, exeFileExt)
	}

	logPwd := filepath.Join(pwd, "logs")

	_, err := os.Stat(logPwd)
	if os.IsNotExist(err) {
		err = os.Mkdir(logPwd, os.ModePerm)
		if err != nil {
			return errors.New("create logs directory failed, " + err.Error())
		}
	}

	loglevel := beego.AppConfig.String("loglevel")
	loglevelNum := 6
	if len(loglevel) != 0 {
		if loglevel == "debug" {
			logs.SetLevel(7)
			loglevelNum = 7
		} else if loglevel == "info" {
			logs.SetLevel(6)
			loglevelNum = 6
		} else if loglevel == "notice" {
			logs.SetLevel(5)
			loglevelNum = 5
		} else if loglevel == "warn" {
			logs.SetLevel(4)
			loglevelNum = 4
		} else if loglevel == "error" {
			logs.SetLevel(3)
			loglevelNum = 3
		} else if loglevel == "critical" {
			logs.SetLevel(2)
			loglevelNum = 2
		} else if loglevel == "alert" {
			logs.SetLevel(1)
			loglevelNum = 1
		} else if loglevel == "emergency" {
			logs.SetLevel(0)
			loglevelNum = 0
		} else {
			logs.Error("invalid log level parameter, option: debug,info,notice,warn,error,critical,alert,emergency")
		}
	}

	log := logs.NewLogger(1000)
	logFilePath := filepath.Join(logPwd, exeFileName+".log")

	logFileConfigMap := make(map[string]interface{})
	logFileConfigMap["filename"] = logFilePath
	logFileConfigMap["level"] = loglevelNum
	logFileConfigMap["daily"] = true
	logFileConfigMap["maxdays"] = 30
	logFileConfigMap["rotate"] = true

	logFileConfigBytes, _ := util.JsonIter.Marshal(logFileConfigMap)
	log.SetLogger("console", `{"level":`+strconv.Itoa(loglevelNum)+`}`)
	logs.SetLogger(logs.AdapterFile, string(logFileConfigBytes))

	return nil

}
