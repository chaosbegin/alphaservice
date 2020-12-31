package impls

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func ArrayJoin(array []interface{}) string {
	result := ""
	for _, v := range array {
		result += fmt.Sprint(v) + ","
	}
	rLen := len(result)
	if rLen > 0 {
		result = result[:rLen-1]
	}

	return result
}

func IntArrayJoin(array []int) string {
	result := ""
	for _, v := range array {
		result += fmt.Sprint(v) + ","
	}
	rLen := len(result)
	if rLen > 0 {
		result = result[:rLen-1]
	}

	return result
}

func ConvertArrayToArrayObject(array []interface{}) []map[string]interface{} {
	arrayObj := make([]map[string]interface{}, 0)
	for _, v := range array {
		//logs.Error("v type:",reflect.TypeOf(v)," value:",v)
		switch v.(type) {
		case map[string]interface{}:
			arrayObj = append(arrayObj, v.(map[string]interface{}))
		}
	}

	return arrayObj
}

func ReconnRandTimeMs() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(5000)
}

func ReConnSleep() {
	time.Sleep(time.Duration(ReconnRandTimeMs()) * time.Millisecond)
}

func HttpReadBody(resp *http.Response) string {
	if resp.Body == nil {
		return ""
	}
	defer resp.Body.Close()
	out, _ := ioutil.ReadAll(resp.Body)
	return string(out)
}
