package httpModel

import (
	"database/sql/driver"
	"demoApp/server/utils"
	"strconv"
	"time"
)

// 时间 字符串转时间戳
type tString time.Time

// scan 结果
func (t *tString) Scan(value interface{}) error {

	*t = tString(value.(time.Time))

	return nil
}

// find 结果
func (t tString) Value() (driver.Value, error) {

	m := time.Time(t)
	return m, nil
}

func TStringFormat(t time.Time) tString {
	return tString(t)
}

// 自定 json 解析结果
func (t tString) MarshalJSON() ([]byte, error) {

	m := time.Time(t)

	return []byte(strconv.Itoa(int(m.Unix()))), nil
}


// 时间转为固定字符串(简历上的开启和结束时间)
type resumeTimeString string

// scan 结果
func (t *resumeTimeString) Scan(value interface{}) error {

	*t = resumeTimeString((value.(time.Time)).Format(utils.RESUME_TIME_FORMAT))

	return nil
}

// find 结果
func (t resumeTimeString) Value() (driver.Value, error) {

	m , err := time.Parse(utils.RESUME_TIME_FORMAT,string(t))
	return m, err
}



// 自定 json 解析结果
//func (t resumeTimeString) MarshalJSON() ([]byte, error) {
//
//	m := time.Time(t).Format(utils.RESUME_TIME_FORMAT)
//
//	return []byte(m), nil
//}



type HttpResultModel struct {
	Result string `json:"result"`
}
