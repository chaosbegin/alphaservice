package main

import (
	_ "alphawolf.com/alphaservice/routers"
	"alphawolf.com/alphaservice/service"
	"alphawolf.com/alphaservice/setup"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
)

var Version string = "8.8.8"
var CommitId string = ""
var BuildTime string = ""
var OS string = ""
var Hardware string = ""

//bee generate appcode -driver=mysql -conn="eyeits:eyeits#1234@tcp(localhost:3306)/eyeits" -level=2 -tables="target"

//../alphabee/alphabee -tpls=../alphabee/tpls generate appcode -driver=mysql -conn="mon:Alphamon#1234@tcp(localhost:3308)/alphamon" -level=2 -tables="target_type"

//export CPLUS_INCLUDE_PATH=/usr//local/lib/node_modules/node-sass/src/libsass/include:$CPLUS_INCLUDE_PATH

func main() {
	err := setup.InitLogger()
	if err != nil {
		fmt.Println("setup logger failed, " + err.Error())
	}

	exePath, _ := os.Executable()
	pwd := filepath.Dir(exePath)
	err = os.Chdir(pwd)
	if err != nil {
		logs.Error("change working dir to ", pwd, " failed, "+err.Error())
	}

	svcConfig := &service.Config{
		Name:        "AlphaService",
		DisplayName: "AlphaService",
		Description: "Alpha monitor web service",
	}

	prg := &setup.Program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logs.Error(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err = s.Install()
			if err != nil {
				logs.Error("service install failed, ", err.Error())
			} else {
				err = setup.SysConfigInstall()
				if err != nil {
					logs.Error("service install setup sysconfig failed, ", err.Error())
				} else {
					logs.Info("service install successful.")
				}
			}
		case "start":
			err = s.Start()
			if err != nil {
				logs.Error("service start failed, ", err.Error())
			} else {
				logs.Info("service start successful.")
			}
			break
		case "stop":
			err = s.Stop()
			if err != nil {
				logs.Error("service stop failed, ", err.Error())
			} else {
				logs.Info("service stop successful.")
			}

			break
		case "restart":
			err = s.Stop()
			if err != nil {
				logs.Error("service stop failed, ", err.Error())
			} else {
				logs.Info("service stop successful.")
			}

			err = s.Start()
			if err != nil {
				logs.Error("service start failed, ", err.Error())
			} else {
				logs.Info("service start successful.")
			}
			break
		case "uninstall":
			err = s.Stop()
			if err != nil {
				logs.Error("service stop failed, ", err.Error())
			} else {
				logs.Info("service stop successful.")
			}
			err = s.Uninstall()
			if err != nil {
				logs.Error("service uninstall failed, ", err.Error())
			} else {
				setup.SysConfigUninstall()
				logs.Info("service uninstall successful.")
			}
			break
		case "version":
			fmt.Printf("AlphaService version: %s-%s %s.%s build:%s\n", OS, Hardware, Version, CommitId, BuildTime)
			return
		case "help":
			PrintHelp()
			return
		//case "id":
		//	id, err := client.HostId("AlphaService")
		//	if err != nil {
		//		logs.Error("generate id failed, ", err.Error())
		//	} else {
		//		fmt.Println("machine id: " + id)
		//	}
		default:
			fmt.Println("option provided but not defined: " + os.Args[1])
			PrintHelp()
		}
		return
	}

	err = s.Run()
	if err != nil {
		logs.Error(err)
	}

}

func PrintHelp() {
	Help := `
USAGE
    alphaservice command

AVAILABLE COMMANDS

    version     Prints the current AlphaService version
    install     Install AlphaService as service
    uninstall   Uninstall AlphaService service
    start       Start AlphaService service 
    stop        Stop AlphaService service 
    restart     Restart AlphaService service 
    help        Prints help information
`

	fmt.Println(Help)
}
