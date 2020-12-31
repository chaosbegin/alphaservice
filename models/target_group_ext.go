package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

func IntArrayJoin(array []int) string {
	result := ""
	for _, v := range array {
		result += fmt.Sprint(v) + ","
	}
	rLen := len(result)
	if rLen > 0 {
		result = result[:rLen-1]
	}

	return result
}

func GetParentTargetGroupIds(o orm.Ormer, targetGroups []TargetGroup) ([]int, error) {
	resIds := make([]int, 0)
	ids := make([]int, 0)
	for _, g := range targetGroups {
		if g.Pid > 0 {
			ids = append(ids, g.Pid)
		}
	}

	tTargetGroups := make([]TargetGroup, 0)

	if len(ids) > 0 {
		_, err := o.Raw("select * from target_group where id in (" + IntArrayJoin(ids) + ")").QueryRows(&tTargetGroups)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}

		if len(tTargetGroups) > 0 {
			for _, g := range tTargetGroups {
				resIds = append(resIds, g.Id)
			}

			tResIds, err := GetParentTargetGroupIds(o, tTargetGroups)
			if err != nil {
				return nil, err
			}

			resIds = append(resIds, tResIds...)
		}
	}

	return resIds, nil
}
