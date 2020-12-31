package impls

import (
	"alphawolf.com/alphaservice/models"
	"errors"
)

var GlobalConfig Config

type Config struct {
}

func (this *Config) GetMasterApiAddr() (string, error) {
	sc, err := models.GetSysConfigById(16)
	if err != nil {
		return "", errors.New("get master api addr failed, " + err.Error())
	}
	if len(sc.Content) < 1 {
		return "", errors.New("master is not available")
	} else {
		return sc.Content, nil
	}
}
