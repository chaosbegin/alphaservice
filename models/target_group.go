package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type TargetGroup struct {
	Id               int       `orm:"column(id);auto"`
	Pid              int       `orm:"column(pid);null"`
	UserId           int       `orm:"column(user_id);null"`
	AutoScan         int       `orm:"column(auto_scan);null"`
	AutoCategory     int       `orm:"column(auto_category);null"`
	CategoryGroupIds string    `orm:"column(category_group_ids);size(512);null"`
	Name             string    `orm:"column(name);size(512);null"`
	StartTime        time.Time `orm:"column(start_time);type(datetime);null"`
	EndTime          time.Time `orm:"column(end_time);type(datetime);null"`
	LastTime         time.Time `orm:"column(last_time);type(datetime);null"`
	HolidayId        int       `orm:"column(holiday_id);null"`
	TimeWindow       int       `orm:"column(time_window);null"`
	ThresholdRewrite int       `orm:"column(threshold_rewrite);null"`
	TargetMatched    int       `orm:"column(target_matched);null"`
	ScanBulkSize     int       `orm:"column(scan_bulk_size);null"`
	IpUnique         int       `orm:"column(ip_unique);null"`
	ScanRange        string    `orm:"column(scan_range);null"`
	ScanSchedule     string    `orm:"column(scan_schedule);null"`
	CategoryRange    string    `orm:"column(category_range);null"`
	ApiId            string    `orm:"column(api_id);null"`
}

func (t *TargetGroup) TableName() string {
	return "target_group"
}

func init() {
	orm.RegisterModel(new(TargetGroup))
}

// AddTargetGroup insert a new TargetGroup into database and returns
// last inserted Id on success.
func AddTargetGroup(m *TargetGroup) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTargetGroupById retrieves TargetGroup by Id. Returns error if
// Id doesn't exist
func GetTargetGroupById(id int) (v *TargetGroup, err error) {
	o := orm.NewOrm()
	v = &TargetGroup{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTargetGroup retrieves all TargetGroup matches certain condition. Returns empty list if
// no records exist
func GetAllTargetGroup(userId int, roleId int, query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(TargetGroup))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else if strings.Contains(k, "__in") {
			vs := strings.Split(v, ",")
			vo := make([]interface{}, len(vs))
			for sk, sv := range vs {
				vo[sk] = sv
			}
			qs = qs.Filter(k, vo...)
		} else {
			qs = qs.Filter(k, v)
		}
	}

	if roleId > 2 {
		targetGroupIds := make([]int, 0)
		targetGroups := make([]TargetGroup, 0)
		_, err = o.Raw("select * from target_group where id in (select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?))", userId).QueryRows(&targetGroups)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}

		if len(targetGroups) < 1 {
			return ml, nil
		}

		for _, g := range targetGroups {
			targetGroupIds = append(targetGroupIds, g.Id)
		}
		to := orm.NewOrm()
		tids, err := GetParentTargetGroupIds(to, targetGroups)
		if err != nil {
			return nil, errors.New("get parent target group failed, " + err.Error())
		}

		targetGroupIds = append(targetGroupIds, tids...)

		logs.Error("targetGroupIds:\n", targetGroupIds)

		if len(targetGroupIds) == 0 {
			return ml, nil
		}

		qs = qs.Filter("id__in", targetGroupIds)
	}

	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	l := make([]TargetGroup, 0)
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateTargetGroup updates TargetGroup by Id and returns error if
// the record to be updated doesn't exist
func UpdateTargetGroupById(m *TargetGroup) (err error) {
	o := orm.NewOrm()
	v := TargetGroup{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTargetGroup deletes TargetGroup by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTargetGroup(id int) (err error) {
	o := orm.NewOrm()
	v := TargetGroup{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TargetGroup{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
