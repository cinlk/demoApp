package utils

import (
	"bytes"
	"goframework/utils"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
)

//

const (
	SEED_STR = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SUCCESS  = "SUCCESS"
	FAILED   = "FAILED"
	RESUME_TIME_FORMAT = "2006-01"
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

func Struct2Map(i interface{}, ignore ...string) *map[string]interface{} {
	var res = map[string]interface{}{}
	v := reflect.Indirect(reflect.ValueOf(i))
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()

	// 忽略 ？？
	// 嵌套不遍历
	for i := 0; i < t.NumField(); i++ {

		switch v.Field(i).Kind() {
		case reflect.String:

			if tag := strings.Split(t.Field(i).Tag.Get("json"),",")[0]; tag != "" {
				res[tag] = v.Field(i).String()
			}
		case reflect.Bool:
			if tag := strings.Split(t.Field(i).Tag.Get("json"),",")[0]; tag != "" {
				res[tag] = v.Field(i).Bool()
			}
		case reflect.Float32, reflect.Float64:
			if tag := strings.Split(t.Field(i).Tag.Get("json"),",")[0]; tag != "" {
				res[tag] = v.Field(i).Float()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if tag := strings.Split(t.Field(i).Tag.Get("json"),",")[0]; tag != "" {
				res[tag] = v.Field(i).Int()
			}
		case reflect.Slice:
			if tag := strings.Split(t.Field(i).Tag.Get("json"),",")[0]; tag != "" {
				// 只有字符串切片
				var tmp []string
				if v.Field(i).Type() == reflect.TypeOf([]string{}){
					for k := 0; k < v.Field(i).Len(); k ++{
						tmp = append(tmp, v.Field(i).Index(k).String())
					}
				}
				res[tag] = tmp
			}

		}

	}

	return &res
}


