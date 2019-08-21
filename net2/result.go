package net

type ErrorInfo struct {
	Msg string `json:"msg"`
	Err error `json:"err"`
}
type Result struct {
	IsOK bool `json:"isOk"`
	Data interface{} `json:"data"`
	ErrorInfo *ErrorInfo `json:"errorInfo"`
}
