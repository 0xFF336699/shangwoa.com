package json2

import (
	"fmt"
	"shangwoa.com/utils/stringnumber"
	"strconv"
)

type onParseUnmatch func(m map[string]interface{}, key string, def interface{}, v interface{}, converted bool)
type onParseError func(m map[string]interface{}, key string, def interface{}, err error)
func DynamicParseString(m map[string]interface{}, key, def string, onConverted onParseUnmatch) string {
	res := def
	converted := false
	var v interface{}
	var b bool
	var matched = false
	if v, b = m[key]; b{
		switch v.(type) {
		case string:
			res = v.(string)
			matched = true
			break
		case int, int32, int64, float32, float64:
			res = stringnumber.ToString(v)
			converted = true
			break
		}
	}
	if !matched{
		onConverted(m, key, def, v, converted)
	}
	return res
}

func DynamicParseBool(m map[string]interface{}, key string, def bool, onConverted onParseUnmatch) bool {
	var res = def
	converted := false
	var v interface{}
	var b bool
	var matched = false
	if v, b = m[key]; b{
		switch v.(type) {
		case bool:
			res = v.(bool)
			matched = true
			break
		default:
			var boolValues = []string{ "1", "TRUE", "true", "True","0","FALSE","false","False", }
			s := fmt.Sprintf("%v", v)
			for _, v := range boolValues{
				if s == v{
					switch s{
					case boolValues[0],boolValues[1],boolValues[2],boolValues[3]:
						res = true
						break
					default:
						res = false
					}
					converted = true
				}
			}
		}
	}
	if !matched{
		onConverted(m, key, def, v, converted)
	}
	return res
}

func DynamicParseInt(m map[string]interface{}, key string, def int, onConverted onParseUnmatch, onError onParseError) int {
	var res = def
	converted := false
	var err error
	var v interface{}
	var b bool
	var matched = false
	if v, b = m[key]; b{
		switch v.(type) {
		case int:
			if d, ok := v.(int);ok{
				res = d
				matched = true
			}
			break
		case float32, float64:
			i := DynamicParseFloat64(m, key, float64(def), onConverted, onError)
			s := fmt.Sprintf("%.0f", i)
			if d, err :=strconv.Atoi(s); err == nil{
				res = d
				matched = true
			}
			break
		case string:
			var i int
			i, err = strconv.Atoi(v.(string))
			if err == nil{
				res = i
				converted = true
			}
			break
		}
	}

	if err != nil && onError != nil{
		onError(m, key, def, err)
	}else if !matched && onConverted != nil{
		onConverted(m, key, def, v, converted)
	}
	return res
}

func DynamicParseInt64(m map[string]interface{}, key string, def int64, onConverted onParseUnmatch, onError onParseError) int64 {
	var res = def
	converted := false
	var err error
	var v interface{}
	var b bool
	matched := false
	if v, b = m[key]; b{
		switch v.(type) {
		case int64:
			res = v.(int64)
			matched = true
			break
		case int,float64,float32:
			if i, ok := v.(int64); ok{
				res = i
				converted = true
			}
			break
		case string:
			var i int64
			i, err = strconv.ParseInt(v.(string), 10, 64)
			if err == nil{
				res = i
				converted = true
			}
			break
		}
	}
	if err != nil && onError != nil{
		onError(m, key, def, err)
	}else if !matched && onConverted != nil{
		onConverted(m, key, def, v, converted)
	}
	return res
}


func DynamicParseFloat64(m map[string]interface{}, key string, def float64, onConverted onParseUnmatch, onError onParseError) float64 {
	var res = def
	converted := false
	var err error
	var v interface{}
	var b bool
	matched := false
	if v, b = m[key]; b{
		switch v.(type) {
		case int,float64,float32:
			res = v.(float64)
			matched = true
			break

		case string:
			var i float64
			i, err = strconv.ParseFloat(v.(string), 10)
			if err == nil{
				res = i
				converted = true
			}
			break
		}
	}

	if err != nil && onError != nil{
		onError(m, key, def, err)
	}else if !matched && onConverted != nil{
		onConverted(m, key, def, v, converted)
	}
	return res
}

func DynamicParseMap(m map[string]interface{}, key string, def map[string]interface{}, onConverted onParseUnmatch) map[string]interface{} {
	var res = def
	var v interface{}
	var b bool
	matched := false
	if v, b = m[key]; b{
		switch v.(type) {
		case map[string]interface{}:
			res = v.(map[string]interface{})
			matched = true
			break
		}
	}
	if !matched && onConverted != nil{
		onConverted(m, key, def, v, true)
	}
	return res
}
