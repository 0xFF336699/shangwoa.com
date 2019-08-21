package net2

import (
	"encoding/json"
	"net/http"
)

type ErrorInfo struct {
	Msg string `json:"msg"`
	Err error `json:"err"`
}
type DataInfo struct {
	Type string `json:"type"`
	Data interface{}
}
type Result struct {
	IsOK bool `json:"isOk"`
	DataInfo *DataInfo `json:"dataInfo"`
	ErrorInfo *ErrorInfo `json:"errorInfo"`
}

func ResponseWrite(err error, data interface{}, w http.ResponseWriter) bool {
	res := &Result{}
	if err !=nil{
		e := &ErrorInfo{Msg:err.Error(), Err:err}
		res.ErrorInfo = e
		res.IsOK = false
	}else {
		res.IsOK = true
	}
	if d, ok := data.(*DataInfo); ok{
		res.DataInfo = d
	}else{
		res.DataInfo = &DataInfo{Data:data}
	}
	b, e := json.Marshal(res)
	if e != nil{
		w.Write([]byte(`{"isOK":false,"errorInfo":{"msg":"parse json error"}`));
	}else{
		w.Write(b)
		return true
	}

	return false
}