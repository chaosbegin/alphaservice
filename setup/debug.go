package setup

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func InitDebug() {
	exePath, _ := os.Executable()
	pwd := filepath.Dir(exePath)
	logPwd := filepath.Join(pwd, "logs")
	if ok, _ := beego.AppConfig.Bool("debug_out"); ok {
		debugFile, _ := os.OpenFile(logPwd+"/debug_out_"+time.Now().Format("2006-1-2_15:04:05")+".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
		syscall.Dup2(int(debugFile.Fd()), 1)
		syscall.Dup2(int(debugFile.Fd()), 2)
	}

	//debug server
	ok, _ := beego.AppConfig.Bool("debug")
	if ok {
		go func() {
			debugAddr := beego.AppConfig.String("debug_addr")
			if len(debugAddr) < 1 {
				debugAddr = "0.0.0.0:6060"
			}
			logs.Info("start debug server on ", debugAddr, " ...")
			err := http.ListenAndServe(debugAddr, nil)
			if err != nil {
				logs.Error("start debug server failed, " + err.Error())
			}
		}()
	}
}
