package impls

import (
	"alphawolf.com/alpha/util"
	"alphawolf.com/alphaservice/models"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strconv"
	"time"
)

type NumRange struct {
	Start int
	End   int
}

type ScheduleItem struct {
	Cycle     int
	FixedTime []*time.Time
	Months    []int
	Weeks     []int
	Days      []int
}

type TimeInterval struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

type MonthDay struct {
	Month int
	Day   int
}

var OperateAclSrv OperateAcl

type OperateAcl struct {
}

func (this *OperateAcl) ConnectCheck(optionId int, targetId int, targetGroupId int, userId int, userGroupIds []int) (bool, error) {
	acls := make([]models.OperateAcl, 0)
	o := orm.NewOrm()
	if len(userGroupIds) > 0 {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? or user_group_id in ("+IntArrayJoin(userGroupIds)+")) and limit_type_id = 1 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return false, errors.New("query acl table failed, " + err.Error())
		}
	} else {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? ) and limit_type_id = 1 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return false, errors.New("query acl table failed, " + err.Error())
		}
	}

	for _, acl := range acls {
		switch acl.OperateTypeId {
		case 1: //deny operate cmd
			ok, err := this.scheduleCheck(acl)
			if err != nil {
				return false, err
			}
			return ok, nil
		default:
			return false, errors.New("not support cmd input operateTypeId:" + strconv.Itoa(acl.OperateTypeId))
		}
	}

	return false, nil
}

func (this *OperateAcl) CmdInputCheck(cmd string, optionId int, targetId int, targetGroupId int, userId int, userGroupIds []int) (bool, error) {
	acls := make([]models.OperateAcl, 0)
	o := orm.NewOrm()
	if len(userGroupIds) > 0 {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? or user_group_id in ("+IntArrayJoin(userGroupIds)+")) and limit_type_id = 2 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return false, errors.New("query acl table failed, " + err.Error())
		}
	} else {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? ) and limit_type_id = 2 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return false, errors.New("query acl table failed, " + err.Error())
		}
	}

	for _, acl := range acls {
		ok, err := this.scheduleCheck(acl)
		if err != nil {
			return false, err
		} else if ok {
			return true, nil
		}

		switch acl.OperateTypeId {
		case 1, 2: //deny operate cmd
			regx, err := regexp.Compile(acl.Pattern)
			if err != nil {
				return false, errors.New("invalid acl pattern, aclId:" + strconv.Itoa(acl.Id))
			}
			if regx.Match([]byte(cmd)) {
				return false, nil
			}
		default:
			return false, errors.New("not support cmd input operateTypeId:" + strconv.Itoa(acl.OperateTypeId))
		}
	}

	return true, nil
}

func (this *OperateAcl) CmdOutputCheck(cmd string, optionId int, targetId int, targetGroupId int, userId int, userGroupIds []int) (string, error) {
	acls := make([]models.OperateAcl, 0)
	o := orm.NewOrm()
	if len(userGroupIds) > 0 {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? or user_group_id in ("+IntArrayJoin(userGroupIds)+")) and limit_type_id = 3 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return "", errors.New("query acl table failed, " + err.Error())
		}
	} else {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? ) and limit_type_id = 3 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return "", errors.New("query acl table failed, " + err.Error())
		}
	}

	for _, acl := range acls {
		ok, err := this.scheduleCheck(acl)
		if err != nil {
			return "", err
		} else if ok {
			return cmd, nil
		}

		switch acl.OperateTypeId {
		case 1: //deny
			regx, err := regexp.Compile(acl.Pattern)
			if err != nil {
				return "", errors.New("invalid acl pattern, aclId:" + strconv.Itoa(acl.Id))
			}
			if regx.Match([]byte(cmd)) {
				return "", nil
			}
		case 2: //delete
			regx, err := regexp.Compile(acl.Pattern)
			if err != nil {
				return "", errors.New("invalid acl pattern, aclId:" + strconv.Itoa(acl.Id))
			}
			cmd = regx.ReplaceAllString(cmd, "")
		default:
			return "", errors.New("not support cmd input operateTypeId:" + strconv.Itoa(acl.OperateTypeId))
		}
	}

	return cmd, nil
}

func (this *OperateAcl) SqlInputCheck(sql string, optionId int, targetId int, targetGroupId int, userId int, userGroupIds []int) (bool, error) {
	acls := make([]models.OperateAcl, 0)
	o := orm.NewOrm()
	if len(userGroupIds) > 0 {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? or user_group_id in ("+IntArrayJoin(userGroupIds)+")) and limit_type_id = 4 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return false, errors.New("query acl table failed, " + err.Error())
		}
	} else {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? ) and limit_type_id = 4 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return false, errors.New("query acl table failed, " + err.Error())
		}
	}

	for _, acl := range acls {
		ok, err := this.scheduleCheck(acl)
		if err != nil {
			return false, err
		} else if ok {
			return true, nil
		}

		switch acl.OperateTypeId {
		case 1, 2: //deny operate cmd
			regx, err := regexp.Compile(acl.Pattern)
			if err != nil {
				return false, errors.New("invalid acl pattern, aclId:" + strconv.Itoa(acl.Id))
			}
			if regx.Match([]byte(sql)) {
				return false, nil
			}
		default:
			return false, errors.New("not support cmd input operateTypeId:" + strconv.Itoa(acl.OperateTypeId))
		}
	}

	return true, nil
}

func (this *OperateAcl) SqlOutputCheck(data []map[string]interface{}, optionId int, targetId int, targetGroupId int, userId int, userGroupIds []int) ([]map[string]interface{}, error) {
	acls := make([]models.OperateAcl, 0)
	o := orm.NewOrm()
	if len(userGroupIds) > 0 {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? or user_group_id in ("+IntArrayJoin(userGroupIds)+")) and limit_type_id = 5 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return nil, errors.New("query acl table failed, " + err.Error())
		}
	} else {
		_, err := o.Raw("select * from operate_acl where (option_id = -1 or option_id = ? or target_id  = ? or  user_id = ? or target_group_id = ? ) and limit_type_id = 5 order by priority", optionId, targetId, userId, targetGroupId).QueryRows(&acls)
		if err != nil {
			return nil, errors.New("query acl table failed, " + err.Error())
		}
	}

	for _, acl := range acls {
		logs.Trace("acl:", acl)
		ok, err := this.scheduleCheck(acl)
		if err != nil {
			return nil, err
		} else if ok {
			return data, nil
		}

		switch acl.OperateTypeId {
		case 1: //deny
			regx, err := regexp.Compile(acl.Pattern)
			if err != nil {
				return nil, errors.New("invalid acl pattern, aclId:" + strconv.Itoa(acl.Id))
			}
			switch acl.OutputMatchTypeId {
			case 1: //text
				dataBytes, err := util.JsonIter.Marshal(data)
				if err != nil {
					return nil, errors.New("marshal data to json string failed, " + err.Error())
				}

				if regx.Match(dataBytes) {
					return nil, errors.New("拒绝访问敏感数据！")
				}

			case 2: //row
				for _, row := range data {
					for _, v := range row {
						if regx.Match([]byte(fmt.Sprint(v))) {
							return nil, errors.New("拒绝访问敏感数据！")
						}
					}
				}
			case 3: //col
				for _, row := range data {
					for k, _ := range row {
						if regx.Match([]byte(k)) {
							return nil, errors.New("拒绝访问敏感数据！")
						}
					}
				}
			default:
				return nil, errors.New("invalid output match type id:" + strconv.Itoa(acl.OutputMatchTypeId))
			}
		case 2: //delete
			regx, err := regexp.Compile(acl.Pattern)
			if err != nil {
				return nil, errors.New("invalid acl pattern, aclId:" + strconv.Itoa(acl.Id))
			}
			switch acl.OutputMatchTypeId {
			case 1, 2: //row
				tdata := make([]map[string]interface{}, 0)
				for _, row := range data {
					canAdd := true
					for _, v := range row {
						if regx.Match([]byte(fmt.Sprint(v))) {
							canAdd = false
							break
						}
					}
					if canAdd {
						tdata = append(tdata, row)
					}
				}
				data = tdata
			case 3: //col
				tdata := make([]map[string]interface{}, 0)
				for _, row := range data {
					trow := make(map[string]interface{})
					for k, v := range row {
						if !regx.Match([]byte(k)) {
							trow[k] = v
						}
					}
					tdata = append(tdata, trow)
				}
				data = tdata
			default:
				return nil, errors.New("invalid output match type id:" + strconv.Itoa(acl.OutputMatchTypeId))
			}
		default:
			return nil, errors.New("invalid cmd input operateTypeId:" + strconv.Itoa(acl.OperateTypeId))
		}
	}

	return data, nil
}

func (this *OperateAcl) scheduleCheck(acl models.OperateAcl) (bool, error) {
	timeNow := time.Now().Local()
	_, monthM, dayNow := timeNow.Date()
	monthNow := int(monthM)

	dayOfWeek := int(time.Now().Weekday())

	nowTotalSec := timeNow.Hour()*3600 + timeNow.Minute()*60 + timeNow.Second()

	if acl.HolidayId > 0 {
		//holiday
		holidays := make([]*MonthDay, 0)
		if acl.HolidayId > -1 {
			holiday, err := models.GetHolidayById(acl.HolidayId)
			if err != nil {
				if err != orm.ErrNoRows {
					errMsg := "acl:" + strconv.Itoa(acl.Id) + " get holiday failed, " + err.Error()
					return false, errors.New(errMsg)
				}

			}

			if holiday != nil && len(holiday.Days) > 0 {
				err = util.JsonIter.Unmarshal([]byte(holiday.Days), &holidays)
				if err != nil {
					errMsg := "acl:" + strconv.Itoa(acl.Id) + " invalid holiday, " + err.Error()
					return false, errors.New(errMsg)
				}
			}

		}

		ok := false
		for _, i := range holidays {
			//logs.Trace("Holiday")
			if monthNow == i.Month && dayNow == i.Day {
				ok = true
				break
			}
		}

		if ok {
			return false, nil
		}

	}

	if len(acl.Schedule) > 0 {
		schedItem := ScheduleItem{}
		err := util.JsonIter.Unmarshal([]byte(acl.Schedule), &schedItem)
		if err != nil {
			return false, errors.New("acl:" + strconv.Itoa(acl.Id) + " invalid schedule format, " + err.Error())
		}

		//Months
		if schedItem.Months != nil && len(schedItem.Months) > 0 && !this.contains(schedItem.Months, monthNow) {
			return false, nil
		}

		//weeks
		if schedItem.Weeks != nil && len(schedItem.Weeks) > 0 && !this.contains(schedItem.Weeks, dayOfWeek) {
			return false, nil
		}

		//days
		if schedItem.Days != nil && len(schedItem.Days) > 0 && !this.contains(schedItem.Days, dayNow) {
			return false, nil
		}
	}

	if acl.TimeWindow > 0 {
		timeWindows := make([]*TimeInterval, 0)
		if acl.TimeWindow > 0 {
			timeWindow, err := models.GetTimeWindowById(acl.TimeWindow)
			if err != nil {
				errMsg := "acl:" + strconv.Itoa(acl.Id) + " get timewindow failed, " + err.Error()
				return false, errors.New(errMsg)
			}
			if len(timeWindow.TimeWindow) > 0 {
				err = util.JsonIter.Unmarshal([]byte(timeWindow.TimeWindow), &timeWindows)
				//logs.Error("timeWindow.TimeWindow:",timeWindow.TimeWindow)
				if err != nil {
					errMsg := "acl:" + strconv.Itoa(acl.Id) + " invalid timewindow, " + err.Error()
					return false, errors.New(errMsg)
				}
			}

		}

		ok := false
		for _, i := range timeWindows {
			startTotalSec := i.Start.Hour()*3600 + i.Start.Minute()*60 + i.Start.Second()
			endTotalSec := i.End.Hour()*3600 + i.End.Minute()*60 + i.End.Second()
			if nowTotalSec >= startTotalSec && nowTotalSec <= endTotalSec {
				ok = true
				break
			}

		}

		if ok {
			return true, nil
		}

	}

	return false, nil

}

func (this *OperateAcl) contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func (this *OperateAcl) GetCols(data []map[string]interface{}) []string {
	cols := make([]string, 0)
	if data != nil && len(data) > 0 {
		for k, _ := range data[0] {
			cols = append(cols, k)
		}
	}

	return cols
}
