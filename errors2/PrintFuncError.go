package errors2

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func PrintFuncError(prefix string, err error, skip int, params ...interface{}) (fn string, pc uintptr, file string, line int, ok bool) {
	pc, file, line, ok = runtime.Caller(skip)
	if ok {
		fn = runtime.FuncForPC(pc).Name()
	}
	var msg string
	if err != nil {
		msg = err.Error()
	}
	bs, e := json.Marshal(params)
	if e != nil {
		fmt.Println("print error info", prefix, fn, line, "marshal data err", params, e.Error(), msg)
	}
	fmt.Println("print error info", prefix, fn, line, string(bs), msg)
	return
}

func PrintInfo(params ...interface{}) (bl bool) {
	bs, e := json.Marshal(params)
	if e != nil {
		var fn = ""
		pc, _, line, ok := runtime.Caller(1)
		if ok {
			fn = runtime.FuncForPC(pc).Name()
		}
		fmt.Println("print error info", fn, line, "marshal data err", params)
	}
	fmt.Println(string(bs))
	return true
}
func TryPrintFuncError(prefix string, err error, skip int, params ...interface{}) (hasError bool) {
	hasError = err != nil
	if hasError {
		PrintFuncError(prefix, err, skip + 1, params)
	}
	return
}

func GetFuncLineInfo(err error, skip int) (hasError, ok bool, fn string, line int, file string, pc uintptr) {
	if err == nil{
		return
	}
	hasError = true
	pc, file, line, ok = runtime.Caller(skip)
	if ok {
		fn = runtime.FuncForPC(pc).Name()
	} else {
		return
	}
	return
}

func GetFuncInfo(skip int) ([]interface{}) {
	var fn string
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		fn = runtime.FuncForPC(pc).Name()
	} else {
		return nil
	}
	return []interface{}{fn, line, file, ok, pc}
}


func PrintLineInfo(skip int, prefix string, params ...interface{}) {
	bs, e := json.Marshal(params)
	if e != nil {
		fmt.Println("PrintLineInfo json marshal err", prefix, params, e.Error())
		return
	}
	var fn string
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		fn = runtime.FuncForPC(pc).Name()
	} else {
		fmt.Println("PrintLineInfo runtime.Caller not ok", prefix, string(bs))
		return
	}
	fmt.Println("PrintLineInfo", prefix, pc, file, line, fn, string(bs))
}
