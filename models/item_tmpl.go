package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ItemTmpl struct {
	AutoIgnore          int       `orm:"column(auto_ignore);null"`
	CategoryId          int       `orm:"column(category_id);null"`
	CmdType             int       `orm:"column(cmd_type);null"`
	CmdRunType          int       `orm:"column(cmd_run_type);null"`
	CmdVersion          string    `orm:"column(cmd_version);null"`
	CmdDetach           int       `orm:"column(cmd_detach);null"`
	CmdDir              string    `orm:"column(cmd_dir);null"`
	CmdEnv              string    `orm:"column(cmd_env);null"`
	CmdVariables        string    `orm:"column(cmd_variables);null"`
	Cols                string    `orm:"column(cols);null"`
	Command             string    `orm:"column(command);null"`
	ConfigPath          string    `orm:"column(config_path);null"`
	ConnTimeout         int       `orm:"column(conn_timeout);null"`
	ConnectOptions      string    `orm:"column(connect_options);size(512);null"`
	Complex             string    `orm:"column(complex);size(4000);null"`
	Dbname              string    `orm:"column(dbname);size(256);null"`
	EndTime             time.Time `orm:"column(end_time);type(datetime);null"`
	ExecTimeout         int       `orm:"column(exec_timeout);null"`
	GroupId             int       `orm:"column(group_id);null"`
	HolidayId           int       `orm:"column(holiday_id);null"`
	Host                string    `orm:"column(host);size(512);null"`
	Id                  int       `orm:"column(id);auto"`
	IgnoreExitstatus    int       `orm:"column(ignore_exitstatus);null"`
	IgnoreStderr        int       `orm:"column(ignore_stderr);null"`
	IsDefault           int       `orm:"column(is_default);null"`
	Passive             int       `orm:"column(passive);null"`
	ItemTypeId          int       `orm:"column(item_type_id);null"`
	MergeCol            int       `orm:"column(merge_col);null"`
	Name                string    `orm:"column(name);size(256);null"`
	ParseMode           string    `orm:"column(parse_mode);size(32);null"`
	Password            string    `orm:"column(password);size(512);null"`
	Pattern             string    `orm:"column(pattern);size(4096);null"`
	Port                string    `orm:"column(port);size(512);null"`
	Preprocess          string    `orm:"column(preprocess);size(4096);null"`
	Retry               int       `orm:"column(retry);null"`
	RetryDelay          int       `orm:"column(retry_delay);null"`
	Saveres             int       `orm:"column(saveres);null"`
	Schedule            string    `orm:"column(schedule);size(4096);null"`
	Series              string    `orm:"column(series);size(256);null"`
	ServiceId           int       `orm:"column(service_id);null"`
	ShortConn           int       `orm:"column(short_conn);null"`
	StartTime           time.Time `orm:"column(start_time);type(datetime);null"`
	Status              int       `orm:"column(status);null"`
	StoreType           int       `orm:"column(store_type);null"`
	Tags                string    `orm:"column(tags);null"`
	Threshold           string    `orm:"column(threshold);null"`
	TimeWindow          int       `orm:"column(time_window);null"`
	TmplTypeId          int       `orm:"column(tmpl_type_id);null"`
	TypeConvs           string    `orm:"column(type_convs);null"`
	UserId              int       `orm:"column(user_id);null"`
	Username            string    `orm:"column(username);size(256);null"`
	Compute             string    `orm:"column(compute);size(256);null"`
	Panel               string    `orm:"column(panel);null"`
	NoticeUserIds       string    `orm:"column(notice_user_ids);null"`
	IgnoreDefaultNotice int       `orm:"column(ignore_default_notice);null"`
}

func (t *ItemTmpl) TableName() string {
	return "item_tmpl"
}

func init() {
	orm.RegisterModel(new(ItemTmpl))
}

// AddItemTmpl insert a new ItemTmpl into database and returns
// last inserted Id on success.
func AddItemTmpl(m *ItemTmpl) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetItemTmplById retrieves ItemTmpl by Id. Returns error if
// Id doesn't exist
func GetItemTmplById(id int) (v *ItemTmpl, err error) {
	o := orm.NewOrm()
	v = &ItemTmpl{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllItemTmpl retrieves all ItemTmpl matches certain condition. Returns empty list if
// no records exist
func GetAllItemTmpl(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) ([]interface{}, error) {
	var err error
	ml := make([]interface{}, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(new(ItemTmpl))
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

	l := make([]ItemTmpl, 0)
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

// UpdateItemTmpl updates ItemTmpl by Id and returns error if
// the record to be updated doesn't exist
func UpdateItemTmplById(m *ItemTmpl) (err error) {
	o := orm.NewOrm()
	v := ItemTmpl{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteItemTmpl deletes ItemTmpl by Id and returns error if
// the record to be deleted doesn't exist
func DeleteItemTmpl(id int) (err error) {
	o := orm.NewOrm()
	v := ItemTmpl{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ItemTmpl{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
