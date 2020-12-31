package setup

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

func isSystemd() bool {
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return true
	}
	if _, err := os.Stat("/proc/1/comm"); err == nil {
		filerc, err := os.Open("/proc/1/comm")
		if err != nil {
			return false
		}
		defer filerc.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(filerc)
		contents := buf.String()

		if strings.Trim(contents, " \r\n") == "systemd" {
			return true
		}
	}
	return false
}

func SysConfigInstall() error {
	if !isSystemd() {
		sysConfigCtx := "ulimit -n 655360\nulimit -u 65536\n"
		err := ioutil.WriteFile(SysConfigPath, []byte(sysConfigCtx), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func SysConfigUninstall() {
	os.Remove(SysConfigPath)
}

const SysConfigPath string = `/etc/sysconfig/AlphaMon`
