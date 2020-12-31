package impls

import (
	"alphawolf.com/alphaservice/models"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"math/big"
)

var OperateAuditSrv OperateAudit

type OperateAudit struct {
	enable    bool
	storeType int
}

func (this *OperateAudit) Start() error {
	this.enable, _ = beego.AppConfig.Bool("operate_audit::enabled")
	if this.enable {
		err := StoreSrv.Start()
		if err != nil {
			return err
		}

		this.storeType, err = beego.AppConfig.Int("operate_audit::store_type")
		if err != nil {
			return errors.New("invalid operate_audit::store_type parameter, " + err.Error())
		}

	}

	return nil
}

func (this *OperateAudit) Add(uuid string, operateTypeId int, user *models.User, target *models.Target, targetOption *models.TargetOption, remoteAddress string, data string) {
	var err error
	if !this.enable {
		return
	}

	row := make(map[string]interface{})
	row["uuid"] = uuid
	row["operate_type_id"] = operateTypeId
	row["user_id"] = user.Id
	row["user_name"] = user.Name
	row["target_id"] = target.Id
	row["target_address"] = target.Address
	row["target_option_id"] = targetOption.Id
	row["target_option_username"] = targetOption.Username
	row["remote_address"] = remoteAddress
	row["data"] = data

	rows := make([]map[string]interface{}, 1)
	rows[0] = row

	msg := &StoreMsg{
		Series: "operate_audit",
		Host:   target.Address,
		Tags:   []string{"target_address", "remote_address", "username"},
		Data:   rows,
	}

	bn := big.NewInt(int64(this.storeType))

	if bn.Bit(0) == 1 { //STORE_TIMEDB = 1
		err = TimedbSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	if bn.Bit(1) == 1 { //STORE_ES = 2
		err = EsSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	if bn.Bit(2) == 1 { //STORE_RDBMS = 4
		err = RDBMSSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

	if bn.Bit(3) == 1 { //STORE_KAFKA = 8
		err = KafkaSrv.PutMsg(msg)
		if err != nil {
			logs.Error(err.Error())
		}
	}

}
