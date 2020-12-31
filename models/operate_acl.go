package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type OperateAcl struct {
	HolidayId         int    `orm:"column(holiday_id);null"`
	Id                int    `orm:"column(id);auto"`
	LimitTypeId       int    `orm:"column(limit_type_id);null"`
	OperateTypeId     int    `orm:"column(operate_type_id);null"`
	OptionId          int    `orm:"column(option_id);null"`
	OutputMatchTypeId int    `orm:"column(output_match_type_id);null"`
	Pattern           string `orm:"column(pattern);null"`
	Priority          int    `orm:"column(priority);null"`
	Schedule          string `orm:"column(schedule);null"`
	TargetGroupId     int    `orm:"column(target_group_Id);null"`
	TargetId          int    `orm:"column(target_id);null"`
	TimeWindow        int    `orm:"column(time_window);null"`
	UserGroupId       int    `orm:"column(user_group_id);null"`
	UserId            int    `orm:"column(user_id);null"`
}

func (t *OperateAcl) TableName() string {
	return "operate_acl"
}

func init() {
	orm.RegisterModel(new(OperateAcl))
}

// AddOperateAcl insert a new OperateAcl into database and returns
// last inserted Id on success.
func AddOperateAcl(m *OperateAcl) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOperateAclById retrieves OperateAcl by Id. Returns error if
// Id doesn't exist
func GetOperateAclById(id int) (v *OperateAcl, err error) {
	o := orm.NewOrm()
	v = &OperateAcl{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOperateAcl retrieves all OperateAcl matches certain condition. Returns empty list if
// no records exist
func GetAllOperateAcl(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(OperateAcl))
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

	l := make([]OperateAcl, 0)
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

// UpdateOperateAcl updates OperateAcl by Id and returns error if
// the record to be updated doesn't exist
func UpdateOperateAclById(m *OperateAcl) (err error) {
	o := orm.NewOrm()
	v := OperateAcl{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOperateAcl deletes OperateAcl by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOperateAcl(id int) (err error) {
	o := orm.NewOrm()
	v := OperateAcl{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&OperateAcl{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
