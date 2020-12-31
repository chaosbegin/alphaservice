package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type HundsunTradingCenter struct {
	CertAlgorithm  string `orm:"column(cert_algorithm);size(255);null"`
	ClientId       string `orm:"column(client_id);size(255);null"`
	CommPassword   string `orm:"column(comm_password);size(255);null"`
	Disable        int    `orm:"column(disable)"`
	Id             int    `orm:"column(id);auto"`
	KeyCertFile    string `orm:"column(key_cert_file);size(512);null"`
	KeyCertPwd     string `orm:"column(key_cert_pwd);size(255);null"`
	LicenseEncrypt string `orm:"column(license_encrypt);size(255);null"`
	LicenseFile    string `orm:"column(license_file);size(512);null"`
	Name           string `orm:"column(name);size(256);null"`
	SafeLevel      string `orm:"column(safe_level);size(32);null"`
	TrustCertFile  string `orm:"column(trust_cert_file);size(512);null"`
	TrustCertPwd   string `orm:"column(trust_cert_pwd);size(255);null"`
	UfxPassword    string `orm:"column(ufx_password);size(255);null"`
}

func (t *HundsunTradingCenter) TableName() string {
	return "hundsun_trading_center"
}

func init() {
	orm.RegisterModel(new(HundsunTradingCenter))
}

// AddHundsunTradingCenter insert a new HundsunTradingCenter into database and returns
// last inserted Id on success.
func AddHundsunTradingCenter(m *HundsunTradingCenter) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHundsunTradingCenterById retrieves HundsunTradingCenter by Id. Returns error if
// Id doesn't exist
func GetHundsunTradingCenterById(id int) (v *HundsunTradingCenter, err error) {
	o := orm.NewOrm()
	v = &HundsunTradingCenter{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHundsunTradingCenter retrieves all HundsunTradingCenter matches certain condition. Returns empty list if
// no records exist
func GetAllHundsunTradingCenter(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(HundsunTradingCenter))
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

	l := make([]HundsunTradingCenter, 0)
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

// UpdateHundsunTradingCenter updates HundsunTradingCenter by Id and returns error if
// the record to be updated doesn't exist
func UpdateHundsunTradingCenterById(m *HundsunTradingCenter) (err error) {
	o := orm.NewOrm()
	v := HundsunTradingCenter{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHundsunTradingCenter deletes HundsunTradingCenter by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHundsunTradingCenter(id int) (err error) {
	o := orm.NewOrm()
	v := HundsunTradingCenter{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&HundsunTradingCenter{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
