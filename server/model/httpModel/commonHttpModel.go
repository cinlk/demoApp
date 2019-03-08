package httpModel

import (
	"database/sql/driver"
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

// 自定 json 解析结果
func (t tString) MarshalJSON() ([]byte, error) {

	m := time.Time(t)

	return []byte(strconv.Itoa(int(m.Unix()))), nil
}


