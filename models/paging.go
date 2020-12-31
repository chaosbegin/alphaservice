package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Page struct {
	Total    int64
	PageSize int64
	PageNo   int64
	Rows     []interface{}
}

// GetAllAlert retrieves all Alert matches certain condition. Returns empty list if
// no records exist
func GetAlertByPage(userId int, roleId int, query map[string]string, fields []string, sortby []string, order []string,
	pageNo int64, pageSzie int64) (Page, error) {
	var err error
	page := Page{
		Total:    0,
		PageSize: 0,
		PageNo:   0,
		Rows:     make([]interface{}, 0),
	}

	ml := make([]interface{}, 0)
	page.Rows = ml

	o := orm.NewOrm()

	qs := o.QueryTable(new(Alert))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}

	if roleId > 2 {
		targetGroupIds := make([]int, 0)
		_, err = o.Raw("select distinct(target_group_id) from target_owner where user_group_id in (select group_id from user_owner where user_id = ?)", userId).QueryRows(&targetGroupIds)
		if err != nil {
			return page, err
		}

		if len(targetGroupIds) == 0 {
			return page, nil
		}

		qs = qs.Filter("target_group_id__in", targetGroupIds)
	}

	page.Total, err = qs.Count()
	if err != nil {
		return page, err
	}

	if pageNo*pageSzie > page.Total {
		return page, errors.New("Invalid pageNo")
	}

	page.PageSize = pageSzie
	page.PageNo = pageNo

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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return page, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return page, errors.New("Error: unused 'order' fields")
		}
	}

	l := make([]Alert, 0)

	qs = qs.OrderBy(sortFields...)
	//logs.Trace("pageNo:",pageNo," pageSize:",pageSzie)
	if _, err = qs.Limit(pageSzie, pageNo*pageSzie).All(&l, fields...); err == nil {
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
		page.Rows = ml
		return page, nil
	}
	page.Rows = ml
	return page, err
}

func GetItemTmplByPage(hideDefault int, groupId int, uid int, query map[string]string, fields []string, sortby []string, order []string,
	pageNo int64, pageSzie int64) (page Page, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ItemTmpl))
	ml := make([]interface{}, 0)
	page.Rows = ml
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}

	if hideDefault == 1 {
		type NameRows struct {
			Name string
			Num  int
		}
		nameRows := make([]NameRows, 0)
		_, err = o.Raw("select name,count(*) as num from item_tmpl where group_id = ? group by name having count(*) > 1", groupId).QueryRows(&nameRows)
		if err != nil && err != orm.ErrNoRows {
			return page, err
		}

		if len(nameRows) > 0 {
			hids := make([]int, 0)
			for _, n := range nameRows {
				//logs.Trace("n:",n)
				trows := make([]*ItemTmpl, 0)
				_, err = o.Raw("select * from item_tmpl where group_id = ? and name = ?", groupId, n.Name).QueryRows(&trows)
				if err != nil && err != orm.ErrNoRows {
					return page, err
				}

				foundOwner := false
				for _, v := range trows {
					if v.UserId == uid {
						foundOwner = true
						break
					}
				}

				for _, v := range trows {
					if foundOwner {
						if v.UserId != uid {
							hids = append(hids, v.Id)
						}
					} else {
						if v.UserId != 1 {
							hids = append(hids, v.Id)
						}
					}
				}
			}

			if len(hids) > 0 {
				for _, i := range hids {
					qs = qs.Exclude("id", i)
				}

			}

		}
	}

	page.Total, err = qs.Count()
	if err != nil {
		return page, err
	}

	if pageNo*pageSzie > page.Total {
		return page, errors.New("Invalid pageNo")
	}

	page.PageSize = pageSzie
	page.PageNo = pageNo

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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return page, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return page, errors.New("Error: unused 'order' fields")
		}
	}

	l := make([]ItemTmpl, 0)

	qs = qs.OrderBy(sortFields...)
	//logs.Trace("pageNo:",pageNo," pageSize:",pageSzie)
	if _, err = qs.Limit(pageSzie, pageNo*pageSzie).All(&l, fields...); err == nil {
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
		page.Rows = ml
		return page, nil
	}
	page.Rows = ml
	return page, err
}

func GetItemByPage(query map[string]string, fields []string, sortby []string, order []string,
	pageNo int64, pageSzie int64) (page Page, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Item))
	ml := make([]interface{}, 0)
	page.Rows = ml
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}

	page.Total, err = qs.Count()
	if err != nil {
		return page, err
	}

	if pageNo*pageSzie > page.Total {
		return page, errors.New("Invalid pageNo")
	}

	page.PageSize = pageSzie
	page.PageNo = pageNo

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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return page, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return page, errors.New("Error: unused 'order' fields")
		}
	}

	l := make([]Item, 0)

	qs = qs.OrderBy(sortFields...)
	//logs.Trace("pageNo:",pageNo," pageSize:",pageSzie)
	if _, err = qs.Limit(pageSzie, pageNo*pageSzie).All(&l, fields...); err == nil {
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
		page.Rows = ml
		return page, nil
	}
	page.Rows = ml
	return page, err
}

func GetUserByPage(groupId int, query map[string]string, fields []string, sortby []string, order []string,
	pageNo int64, pageSzie int64) (page Page, err error) {
	ml := make([]interface{}, 0)
	page.Rows = ml
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}

	if groupId != -1024 {
		userIds := make([]int, 0)
		_, err = o.Raw("select user_id from user_owner where group_id = ?", groupId).QueryRows(&userIds)
		if err != nil {
			return page, err
		}
		if len(userIds) < 1 {
			return page, nil
		}
		qs = qs.Filter("id__in", userIds)
	}

	page.Total, err = qs.Count()
	if err != nil {
		return page, err
	}

	if pageNo*pageSzie > page.Total {
		return page, errors.New("Invalid pageNo")
	}

	page.PageSize = pageSzie
	page.PageNo = pageNo

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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return page, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return page, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return page, errors.New("Error: unused 'order' fields")
		}
	}

	l := make([]User, 0)

	qs = qs.OrderBy(sortFields...)
	//logs.Trace("pageNo:",pageNo," pageSize:",pageSzie)
	if _, err = qs.Limit(pageSzie, pageNo*pageSzie).All(&l, fields...); err == nil {
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
		page.Rows = ml
		return page, nil
	}
	page.Rows = ml
	return page, err
}
