package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AlertHistory struct {
	ActionData      string    `orm:"column(action_data);size(1024);null"`
	ConfirmTime     time.Time `orm:"column(confirm_time);type(datetime);null"`
	ConfirmUid      int       `orm:"column(confirm_uid);null"`
	ConfirmUsername string    `orm:"column(confirm_username);size(256);null"`
	Hash            string    `orm:"column(hash);size(128);null"`
	Host            string    `orm:"column(host);size(512);null"`
	Id              int       `orm:"column(id);auto"`
	ItemGroupId     int       `orm:"column(item_group_id);null"`
	ItemGroupName   string    `orm:"column(item_group_name);size(512);null"`
	ItemId          int       `orm:"column(item_id);null"`
	Lasttime        time.Time `orm:"column(lasttime);type(datetime);null"`
	Level           int       `orm:"column(level);null"`
	Message         string    `orm:"column(message);null"`
	Mode            string    `orm:"column(mode);size(64);null"`
	Name            string    `orm:"column(name);size(512);null"`
	RecordTime      time.Time `orm:"column(record_time);type(datetime);null"`
	Responsible     string    `orm:"column(responsible);size(512);null"`
	Series          string    `orm:"column(series);size(1024);null"`
	Starttime       time.Time `orm:"column(starttime);type(datetime);null"`
	AlertType       int       `orm:"column(alert_type);null"`
	ConfigId        int       `orm:"column(config_id);null"`
	Status          int       `orm:"column(status);null"`
	TagColumn       string    `orm:"column(tag_column);size(1024);null"`
	Tagging         string    `orm:"column(tagging);null"`
	TargetGroupId   int       `orm:"column(target_group_id);null"`
	TargetGroupName string    `orm:"column(target_group_name);size(512);null"`
	TargetId        int       `orm:"column(target_id);null"`
	TargetName      string    `orm:"column(target_name);size(512);null"`
	ThresPath       int       `orm:"column(thres_path);null"`
	Times           int       `orm:"column(times);null"`
	UserGroupId     int       `orm:"column(user_group_id);null"`
	UserFlag        string    `orm:"column(user_flag);null"`
}

func (t *AlertHistory) TableName() string {
	return "alert_history"
}

func init() {
	orm.RegisterModel(new(AlertHistory))
}

// AddAlertHistory insert a new AlertHistory into database and returns
// last inserted Id on success.
func AddAlertHistory(m *AlertHistory) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAlertHistoryById retrieves AlertHistory by Id. Returns error if
// Id doesn't exist
func GetAlertHistoryById(id int) (v *AlertHistory, err error) {
	o := orm.NewOrm()
	v = &AlertHistory{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAlertHistory retrieves all AlertHistory matches certain condition. Returns empty list if
// no records exist
func GetAllAlertHistory(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(AlertHistory))
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

	l := make([]AlertHistory, 0)
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

// UpdateAlertHistory updates AlertHistory by Id and returns error if
// the record to be updated doesn't exist
func UpdateAlertHistoryById(m *AlertHistory) (err error) {
	o := orm.NewOrm()
	v := AlertHistory{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAlertHistory deletes AlertHistory by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAlertHistory(id int) (err error) {
	o := orm.NewOrm()
	v := AlertHistory{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AlertHistory{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
