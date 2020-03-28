package iinterface

import (
	"fmt"
	"strconv"
)

func MustToInt(v interface{}, def int) (int) {
	r := def
	switch v.(type) {
	case int:
		r = v.(int)
		break
	case int64:
		if n, ok := v.(int); ok{
			r = n
		}
		break
	case string:
		n, e := strconv.Atoi(v.(string))
		if e == nil{
			r = n
		}
		break
	default:
		fmt.Printf("warning MustToInt got unknow type v is %+v", v)
		break
	}
	return r
}

func MustToBool(v interface{}, def bool) bool  {
	r := def
	switch v.(type) {
	case int:
		r = v.(int) == 1
		break
	case int64:
		r = v.(int64) == 1
		break
	case string:
		if b, e := strconv.ParseBool(v.(string)); e == nil{
			r = b
		}
		break
	case bool:
		r = v.(bool)
		break
	default:
		fmt.Printf("warning MustToBool got unknow type v is %+v", v)
		break
	}
	return r
}

func MustToString(v interface{}, def string) string {
	r := def
	switch v.(type) {
	case int:
		r = strconv.Itoa(v.(int))
		break
	case int64:
		r = strconv.FormatInt(v.(int64), 10)
		break
	case bool:
		r = strconv.FormatBool(v.(bool))
		break
	case string:
		r = v.(string)
		break
	default:
		fmt.Printf("warning MustToString got unknow type v is %+v", v)
		break
	}
	return r
}