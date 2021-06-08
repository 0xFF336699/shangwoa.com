package consul

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"shangwoa.com/http2"
	"shangwoa.com/json2"
	"shangwoa.com/log2"
	"strings"
)

// KVValue consul里实际储存的键值。这个是当前暂时默认使用的结构，如果有其它的结构则需要另行定义。
// 如果是json，则能够储存bool和数字，如果value直接是值，则都会被认为是string类型，所以推荐使用json结构体储存value
// mq现在是用KVValue来配置
type KVValue struct {
	// Int 如果value是int
	Int int
	// Bool 如果value是bool
	Bool bool
	// String 如果value是string
	String string
	// Float64 如果value是float64，看似像是数字都是这个类型。
	Float64 float64 `json:"float_64"`
	Num     int     `json:"num"`
	Bl      bool    `json:"bl"`
	Str     string  `json:"str"`
	Float32 float32 `json:"float_32"`
	// ValueKind 使用时可以先判断一下value的类型，再到Int、Bool、String、Float64里去找
	ValueKind reflect.Kind // reflect.Bool reflect.String ect
	// Value 原始数据，使用Int Bool String Float64就行，这个字段仅供查阅。
	Value interface{} `json:"value"`
	// Description 数据的一些解释，协助理解用
	Description string `json:"description"`
	// Name 字段或者数据名，协助理解用
	Name string `json:"name"`
	// Extra 额外的信息，协助理解用
	Extra string `json:"extra"`
}

type  KVInt struct{
	Name string `json:"name"`
	Explanation string `json:"explanation"`
	Value int `json:"value"`
}

type  KVFloat64 struct{
	Name string `json:"name"`
	Explanation string `json:"explanation"`
	Value float64 `json:"value"`
}

type  KVBool struct{
	Name string `json:"name"`
	Explanation string `json:"explanation"`
	Value bool `json:"value"`
}

type  KVString struct{
	Name string `json:"name"`
	Explanation string `json:"explanation"`
	Value string `json:"value"`
}

type KVRedis struct {
	Addr string `json:"addr"`
	PW   string `json:"pw"`
	DB   int    `json:"db"`
	// Description 数据的一些解释，协助理解用
	Description string `json:"description"`
	// Name 字段或者数据名，协助理解用
	Name string `json:"name"`
	// Extra 额外的信息，协助理解用
	Extra string `json:"extra"`
}

type KV struct {
	Err *error
	Key string
	// KVStruct 自定义结构体，必须是个指针。优先判断使用这个
	KVStruct interface{}
	// KVValue 已提供的默认数据结构体。第二判断使用这个
	KVValue *KVValue
	// KVJSONB jsonb格式数据。 第三判断使用这个
	KVJSONB *json2.JSONB
	// KVStringValue 如果以上KVValue、KVStruct、KVJSONB都不赋值的话，第四则默认使用这个
	KVStringValue string
	KVInt *KVInt
	KVFloat64 *KVFloat64
	KVBool *KVBool
	KVString *KVString
}

type TxnResult struct {
	Errors      []Errors  `json:"errors"`
	Index       int       `json:"Index"`
	LastContact int       `json:"LastContact"`
	KnownLeader bool      `json:"KnownLeader"`
	Results     []Results `json:"Results"`
}

type Results struct {
	KVRes KVRes `json:"KV"`
}

type KVRes struct {
	LockIndex   int    `json:"LockIndex"`
	Key         string `json:"Key"`
	Flags       int    `json:"Flags"`
	Value       string `json:"Value"`
	CreateIndex int    `json:"CreateIndex"`
	ModifyIndex int    `json:"ModifyIndex"`
}

type Errors struct {
	OpIndex int    `json:"OpIndex"`
	What    string `json:"What"`
}


func GetKeys(url string, m map[string]*KV) error {
	l := len(m)
	if l == 0 {
		err := errors.New("not enough paramaters")
		return err
	}
	body := ""
	for k, v := range m {
		v.Key = k
		if len(body) > 0 {
			body += ","
		}
		body += `{
			"KV": {
			  "Verb": "get",
			  "Key": "` + k + `"
			}
		  }`
	}
	body = "[" + body + "]"
	var jsonStr = []byte(body)

	r, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	r.Header.Set("X-Custom-Header", "myvalue")
	r.Header.Set("Content-Type", "application/json")
	b, err := http2.ClientDo(r, http2.GetEmptyClient())
	if err != nil {
		fmt.Println("consul error is", string(b))
		return err
	}
	var txn TxnResult
	err = json.Unmarshal(b, &txn)
	if err != nil {
		fmt.Println("txn error", err)
		return err
	}
	if len(txn.Errors) > 0 {
		info := ""
		for _, e := range txn.Errors {
			info += fmt.Sprintf("txn error index is %d and Wath is %s \n", e.OpIndex, e.What)
		}
		info = "txn error is \n" + info
		err = errors.New(info)
		log2.Error(err)
		return err
	}

	for i := 0; i < l; i++ {
		res := &txn.Results[i].KVRes
		err = parseValue(res, m[res.Key])
		if err != nil {
			return err
		}
	}
	return nil
}

func parseValue2(res *KVRes, kv *KV) (err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(res.Value)
	if err != nil {
		kv.Err = &err
		return
	}
	if kv.KVStruct != nil {
		err = json.Unmarshal(decodeBytes, kv.KVStruct)
		if err != nil {
			kv.Err = &err
			return
		}
	} else if kv.KVJSONB != nil {
		err = json.Unmarshal(decodeBytes, kv.KVJSONB)
		if err != nil {
			kv.Err = &err
			return
		}
	} else if kv.KVValue != nil {
		v := kv.KVValue
		err = json.Unmarshal(decodeBytes, v)
		if err != nil {
			kv.Err = &err
			return
		}
		switch v.Value.(type) {
		case int:
			v.Int = v.Value.(int)
			v.ValueKind = reflect.Int
		case float64:
			v.Float64 = v.Value.(float64)
			v.ValueKind = reflect.Float64
		case string:
			v.String = v.Value.(string)
			v.ValueKind = reflect.String
		case bool:
			v.Bool = v.Value.(bool)
			v.ValueKind = reflect.Bool
		default:
			err = errors.New("unknow type")
			kv.Err = &err
			return
		}
	}else if kv.KVInt != nil{
		err = json.Unmarshal(decodeBytes, kv.KVInt)
		if err != nil {
			kv.Err = &err
			return
		}
	}else {
		kv.KVStringValue = string(decodeBytes)
	}
	return nil
}


func parseValue(res *KVRes, kv *KV) (err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(res.Value)
	if err != nil {
		kv.Err = &err
		return
	}
	var i interface{}
	if kv.KVStruct != nil{
		i = kv.KVStruct
	}else if kv.KVJSONB != nil {
		i = kv.KVJSONB
	}else if kv.KVValue != nil{
		i = kv.KVValue
	}else if kv.KVInt != nil{
		i = kv.KVInt
	}else if kv.KVFloat64 != nil{
		i = kv.KVFloat64
	}else if kv.KVBool != nil{
		i = kv.KVBool
	}else if kv.KVString != nil{
		i = kv.KVString
	}

	if i != nil{
		err = json.Unmarshal(decodeBytes, i)
		if err != nil {
			kv.Err = &err
			return
		}
		if kv.KVValue == nil{
			return
		}
		v := kv.KVValue
		switch v.Value.(type) {
		case int:
			v.Int = v.Value.(int)
			v.ValueKind = reflect.Int
		case float64:
			v.Float64 = v.Value.(float64)
			v.ValueKind = reflect.Float64
		case string:
			v.String = v.Value.(string)
			v.ValueKind = reflect.String
		case bool:
			v.Bool = v.Value.(bool)
			v.ValueKind = reflect.Bool
		default:
			err = errors.New("unknow type")
			kv.Err = &err
			return
		}
	}else{
		kv.KVStringValue = string(decodeBytes)
	}
	return nil
}

func pathToWords(p string) (s string) {
	words := strings.Split(p, "/")
	s = ""
	for _, v := range words{
		s += strings.Title(v)
	}
	return
}