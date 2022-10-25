package stringnumber

import (
	"fmt"
	"reflect"
	"strconv"
)


func ToString(arg interface{}) string {
	var tmp = reflect.Indirect(reflect.ValueOf(arg)).Interface()
	switch v := tmp.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	//case time.Time:
	//	if len(timeFormat) == 1 {
	//		return v.Format(timeFormat[0])
	//	}
	//	return v.Format("2006-01-02 15:04:05")
	//case jsoncrack.Time:
	//	if len(timeFormat) == 1 {
	//		return v.Time().Format(timeFormat[0])
	//	}
	//	return v.Time().Format("2006-01-02 15:04:05")
	case fmt.Stringer:
		return v.String()
	case reflect.Value:
		return ToString(v.Interface())
	default:
		return ""
	}
}