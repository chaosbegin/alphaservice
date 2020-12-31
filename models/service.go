package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Service struct {
	Id            int    `orm:"column(id);auto"`
	ServiceName   string `orm:"column(service_name);size(255);null"`
	ServiceUrl    string `orm:"column(service_url);size(512);null"`
	LogOutput     string `orm:"column(log_output);null"`
	CkDistributed int    `orm:"column(ck_distributed);null"`
	CkDatabase    string `orm:"column(ck_database);null"`
	CkCluster     string `orm:"column(ck_cluster);null"`
	CkEngine      string `orm:"column(ck_engine);null"`
	CkPartitionBy string `orm:"column(ck_partition_by);null"`
	CkOrderBy     string `orm:"column(ck_order_by);null"`
	CkSettings    string `orm:"column(ck_settings);null"`
}

func (t *Service) TableName() string {
	return "service"
}

func init() {
	orm.RegisterModel(new(Service))
}

// AddService insert a new Service into database and returns
// last inserted Id on success.
func AddService(m *Service) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetServiceById retrieves Service by Id. Returns error if
// Id doesn't exist
func GetServiceById(id int) (v *Service, err error) {
	o := orm.NewOrm()
	v = &Service{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllService retrieves all Service matches certain condition. Returns empty list if
// no records exist
func GetAllService(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(Service))
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

	l := make([]Service, 0)
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

// UpdateService updates Service by Id and returns error if
// the record to be updated doesn't exist
func UpdateServiceById(m *Service) (err error) {
	o := orm.NewOrm()
	v := Service{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteService deletes Service by Id and returns error if
// the record to be deleted doesn't exist
func DeleteService(id int) (err error) {
	o := orm.NewOrm()
	v := Service{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Service{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
