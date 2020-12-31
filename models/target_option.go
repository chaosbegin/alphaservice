package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type TargetOption struct {
	Id             int       `orm:"column(id);auto"`
	TypeId         int       `orm:"column(type_id);null"`
	TargetId       int       `orm:"column(target_id);null"`
	UserId         int       `orm:"column(user_id);null"`
	ItemTypeId     int       `orm:"column(item_type_id);null"`
	ServiceId      int       `orm:"column(service_id);null"`
	Isdefault      int       `orm:"column(isdefault);null"`
	AddrType       int       `orm:"column(addr_type);null"`
	Port           string    `orm:"column(port);size(512);null"`
	Username       string    `orm:"column(username);size(512);null"`
	Password       string    `orm:"column(password);size(512);null"`
	Dbname         string    `orm:"column(dbname);size(512);null"`
	ConnectOptions string    `orm:"column(connect_options);null"`
	AutoPwd        int       `orm:"column(auto_pwd);null"`
	Version        string    `orm:"column(version);size(512);null"`
	ChangeTime     time.Time `orm:"column(change_time);type(timestamp);null"`
	Strategy       string    `orm:"column(strategy);null"`
}

func (t *TargetOption) TableName() string {
	return "target_option"
}

func init() {
	orm.RegisterModel(new(TargetOption))
}

// AddTargetOption insert a new TargetOption into database and returns
// last inserted Id on success.
func AddTargetOption(m *TargetOption) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTargetOptionById retrieves TargetOption by Id. Returns error if
// Id doesn't exist
func GetTargetOptionById(id int) (v *TargetOption, err error) {
	o := orm.NewOrm()
	v = &TargetOption{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTargetOption retrieves all TargetOption matches certain condition. Returns empty list if
// no records exist
func GetAllTargetOption(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(TargetOption))
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

	l := make([]TargetOption, 0)
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

// UpdateTargetOption updates TargetOption by Id and returns error if
// the record to be updated doesn't exist
func UpdateTargetOptionById(m *TargetOption) (err error) {
	o := orm.NewOrm()
	v := TargetOption{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTargetOption deletes TargetOption by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTargetOption(id int) (err error) {
	o := orm.NewOrm()
	v := TargetOption{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TargetOption{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
