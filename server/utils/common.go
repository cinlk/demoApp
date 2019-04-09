package utils

import (
	"bytes"
	"goframework/utils"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

//

const (
	SEED_STR = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SUCCESS  = "SUCCESS"
	FAILED   = "FAILED"
)

// 加密密码

func RandStrings(length int) string {

	b := bytes.Buffer{}
	rand.Seed(time.Now().UnixNano())
	lens := len(SEED_STR)
	for i := 0; i < length; i++ {
		b.WriteByte(SEED_STR[rand.Intn(lens)])
	}
	return b.String()

}

func JsonResponse(w http.ResponseWriter, data interface{}, code int) {
	utils.JsonResponse(w, data, code)
}


func Struct2Map(i interface{}) *map[string]interface{}{
	var res = map[string]interface{}{}
	v := reflect.Indirect(reflect.ValueOf(i))
	if v.Kind() != reflect.Struct{
		return  nil
	}
	t := v.Type()

	// 嵌套不遍历
	for  i:=0;   i < t.NumField(); i++{

		switch v.Field(i).Kind() {
		case reflect.String:
			if tag := t.Field(i).Tag.Get("json"); tag != ""{
				res[tag] = v.Field(i).String()
			}
		case reflect.Bool:
			if tag := t.Field(i).Tag.Get("json"); tag != ""{
				res[tag] = v.Field(i).Bool()
			}
		case reflect.Float32, reflect.Float64:
			if tag := t.Field(i).Tag.Get("json"); tag != ""{
				res[tag] = v.Field(i).Float()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if tag := t.Field(i).Tag.Get("json"); tag != ""{
				res[tag] = v.Field(i).Int()
			}
		}

	}

	return &res
}