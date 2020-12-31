package impls

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"errors"
	"github.com/astaxie/beego/orm"
	"strings"
)

type Panel struct {
	Name  string
	Id    string
	Param map[string]interface{}
}

func GetPanelByAlertId(id int) ([]Panel, error) {
	alert, err := models.GetAlertById(id)
	if err != nil {
		return nil, errors.New("invalid alert id, " + err.Error())
	}

	item, _ := models.GetItemById(alert.ItemId)

	if item == nil {
		return GetPanelByTargetId(alert.TargetId)
	} else if len(item.Panel) < 1 {
		return GetPanelByTargetId(item.TargetId)
	}

	target, _ := models.GetTargetById(alert.TargetId)

	ps := make([]Panel, 0)
	err = util.JsonIter.Unmarshal([]byte(item.Panel), &ps)
	if err != nil {
		return nil, errors.New("invalid panel data, " + err.Error())
	} else {
		for i := 0; i < len(ps); i++ {
			if ps[i].Param == nil {
				param := make(map[string]interface{}, 6)
				param["host"] = alert.Host
				param["name"] = alert.Name
				param["item_id"] = alert.ItemId
				param["group"] = alert.TargetGroupName
				if item != nil {
					param["item"] = item.Name
				}
				if target != nil {
					param["target_type"] = strings.Split(target.TargetType, ";")
				}
				ps[i].Param = param
			} else {
				ps[i].Param["host"] = alert.Host
				ps[i].Param["name"] = alert.Name
				ps[i].Param["item"] = item.Name

				if item != nil {
					ps[i].Param["item_id"] = item.Id
				}
				ps[i].Param["group"] = alert.TargetGroupName
				if target != nil {
					ps[i].Param["target_type"] = strings.Split(target.TargetType, ";")
				}
			}
		}
		return ps, nil
	}

}

func GetPanelByItemId(id int) ([]Panel, error) {
	item, err := models.GetItemById(id)
	if err != nil {
		return nil, errors.New("query item failed, " + err.Error())
	}

	if len(item.Panel) < 1 {
		return GetPanelByTargetId(item.TargetId)
	}

	target, err := models.GetTargetById(item.TargetId)
	if err != nil {
		return nil, errors.New("query target by id failed, " + err.Error())
	}

	ps := make([]Panel, 0)
	err = util.JsonIter.Unmarshal([]byte(item.Panel), &ps)
	if err != nil {
		return nil, errors.New("invalid panel data, " + err.Error())
	} else {
		for i := 0; i < len(ps); i++ {
			if ps[i].Param == nil {
				param := make(map[string]interface{}, 6)
				param["host"] = target.Address
				param["name"] = target.Name
				param["item"] = item.Name
				param["item_id"] = item.Id
				param["group"] = GetTargetGroupName(target)
				param["target_type"] = strings.Split(target.TargetType, ";")
				ps[i].Param = param
			} else {
				ps[i].Param["host"] = target.Address
				ps[i].Param["name"] = target.Name
				ps[i].Param["item"] = item.Name
				ps[i].Param["item_id"] = item.Id
				ps[i].Param["group"] = GetTargetGroupName(target)
				ps[i].Param["target_type"] = strings.Split(target.TargetType, ";")
			}
		}
		return ps, nil
	}

}

func GetPanelByTargetId(id int) ([]Panel, error) {
	target, err := models.GetTargetById(id)
	if err != nil {
		return nil, errors.New("get target by id failed, " + err.Error())
	}

	tt := "default"
	if len(target.TargetType) > 1 {
		tt = target.TargetType
	}

	rows := make([]Panel, 0)

	tts := strings.Split(tt, ";")
	o := orm.NewOrm()
	for _, t := range tts {
		pss := make([]string, 0)
		_, err = o.Raw("select panel from target_panel where target_type = ?", t).QueryRows(&pss)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}

		for _, ps := range pss {
			if len(ps) > 0 {
				row := make([]Panel, 0)
				err = util.JsonIter.Unmarshal([]byte(ps), &row)
				if err != nil {
					return nil, errors.New("invalid panel format, " + err.Error())
				}

				for i := 0; i < len(row); i++ {
					if row[i].Param == nil {
						param := make(map[string]interface{}, 4)
						param["host"] = target.Address
						param["name"] = target.Name
						param["group"] = GetTargetGroupName(target)
						param["target_type"] = strings.Split(target.TargetType, ";")
						row[i].Param = param
					} else {
						row[i].Param["host"] = target.Address
						row[i].Param["name"] = target.Name
						row[i].Param["group"] = GetTargetGroupName(target)
						row[i].Param["target_type"] = strings.Split(target.TargetType, ";")

					}
				}

				if len(row) > 0 {
					rows = append(rows, row...)
				}

			}
		}
	}

	if len(rows) < 1 {
		pss := make([]string, 0)
		_, err = o.Raw("select panel from target_panel where target_type = 'default'").QueryRows(&pss)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}

		for _, ps := range pss {
			if len(ps) > 0 {
				row := make([]Panel, 0)
				err = util.JsonIter.Unmarshal([]byte(ps), &row)
				if err != nil {
					return nil, errors.New("invalid panel format, " + err.Error())
				}

				for i := 0; i < len(row); i++ {
					if row[i].Param == nil {
						param := make(map[string]interface{}, 4)
						param["host"] = target.Address
						param["name"] = target.Name
						param["group"] = GetTargetGroupName(target)
						param["target_type"] = strings.Split(target.TargetType, ";")
						row[i].Param = param
					} else {
						row[i].Param["host"] = target.Address
						row[i].Param["name"] = target.Name
						row[i].Param["group"] = GetTargetGroupName(target)
						row[i].Param["target_type"] = strings.Split(target.TargetType, ";")

					}
				}

				if len(row) > 0 {
					rows = append(rows, row...)
				}

			}
		}

	}

	return rows, nil
}

func GetTargetGroupName(target *models.Target) string {
	targetGroup, err := models.GetTargetGroupById(target.GroupId)
	if err != nil {
		return ""
	}
	if targetGroup.Pid > 0 {
		pTargetGroup, err := models.GetTargetGroupById(targetGroup.Pid)
		if err != nil {
			return targetGroup.Name
		} else {
			return pTargetGroup.Name + " -> " + targetGroup.Name
		}
	} else {
		return targetGroup.Name
	}
}
