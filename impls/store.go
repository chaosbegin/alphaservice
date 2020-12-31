package impls

import "github.com/astaxie/beego/logs"

var StoreSrv *Store

func init() {
	StoreSrv = &Store{}
}

type Store struct {
}

func (this *Store) Start() error {
	err := RDBMSSrv.Start()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	err = TimedbSrv.Start()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	err = EsSrv.Start()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	err = KafkaSrv.Start()
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	return nil
}

func (this *Store) AddData(msg *StoreMsg) {
	if RDBMSSrv.enabled {
		err := RDBMSSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	if TimedbSrv.enabled {
		err := TimedbSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	if EsSrv.enabled {
		err := EsSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	if KafkaSrv.enabled {
		err := KafkaSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

}

func (this *Store) LogAdd(level int, msg ...interface{}) {
	switch level {
	case logs.LevelEmergency:
		logs.Emergency(msg)
	case logs.LevelAlert:
		logs.Alert(msg)
	case logs.LevelCritical:
		logs.Critical(msg)
	case logs.LevelError:
		logs.Error(msg)
	case logs.LevelWarn:
		logs.Warn(msg)
	case logs.LevelNotice:
		logs.Notice(msg)
	case logs.LevelInfo:
		logs.Info(msg)
	case logs.LevelDebug:
		logs.Debug(msg)
	default:
		logs.Info(msg)
	}

}
