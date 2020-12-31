package impls

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

//var NotAuthItemtypeMap map[int]int
//var NotAuthItemtypeString string

//func init() {
//	NotAuthItemtypeString = "4,8,11,12,13,50,52"
//	NotAuthItemtypeMap = make(map[int]int)
//	NotAuthItemtypeMap[4] = 0
//	NotAuthItemtypeMap[8] = 0
//	NotAuthItemtypeMap[11] = 0
//	NotAuthItemtypeMap[12] = 0
//	NotAuthItemtypeMap[13] = 0
//	NotAuthItemtypeMap[50] = 0
//	NotAuthItemtypeMap[52] = 0
//}

func ItemTypeNeedAuth(itemTypeId int) bool {
	count := 0
	o := orm.NewOrm()
	err := o.Raw("select count(*) from item_type where no_auth = 1 and id = ?", itemTypeId).QueryRow(&count)
	if err != nil {
		logs.Error("query no auth item type failed, " + err.Error())
		return false
	}

	if count > 0 {
		return false
	} else {
		return true
	}
}
