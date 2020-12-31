package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PwdOption struct {
	ChangeTime time.Time `orm:"column(change_time);type(timestamp);null"`
	ConnOption string    `orm:"column(conn_option);null"`
	Dbname     string    `orm:"column(dbname);null"`
	Expire     int       `orm:"column(expire);null"`
	Id         int       `orm:"column(id);auto"`
	Password   string    `orm:"column(password);null"`
	Port       int       `orm:"column(port);null"`
	ServiceId  int       `orm:"column(service_id);null"`
	Strategy   string    `orm:"column(strategy);null"`
	TargetId   int       `orm:"column(target_id);null"`
	TypeId     int       `orm:"column(type_id);null"`
	Username   string    `orm:"column(username);null"`
	Version    string    `orm:"column(version);size(512)"`
}

func (t *PwdOption) TableName() string {
	return "pwd_option"
}

func init() {
	orm.RegisterModel(new(PwdOption))
}

// AddPwdOption insert a new PwdOption into database and returns
// last inserted Id on success.
func AddPwdOption(m *PwdOption) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPwdOptionById retrieves PwdOption by Id. Returns error if
// Id doesn't exist
func GetPwdOptionById(id int) (v *PwdOption, err error) {
	o := orm.NewOrm()
	v = &PwdOption{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPwdOption retrieves all PwdOption matches certain condition. Returns empty list if
// no records exist
func GetAllPwdOption(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(PwdOption))
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

	l := make([]PwdOption, 0)
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

// UpdatePwdOption updates PwdOption by Id and returns error if
// the record to be updated doesn't exist
func UpdatePwdOptionById(m *PwdOption) (err error) {
	o := orm.NewOrm()
	v := PwdOption{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePwdOption deletes PwdOption by Id and returns error if
// the record to be deleted doesn't exist
func DeletePwdOption(id int) (err error) {
	o := orm.NewOrm()
	v := PwdOption{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PwdOption{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
