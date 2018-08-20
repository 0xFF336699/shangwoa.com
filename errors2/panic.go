package errors2

import (
	"runtime"
	"shangwoa.com/log2"
)
// PanicError 一般是初始化的时候遇到炒鸡错误需要跳出，初始化完成后不应该使用这个方法跳出
func PanicError(err error) {
	_, file, line, _ := runtime.Caller(2)
	log2.Printf("panic error is %#v \n file is %s \n line is %d \n", err, file, line)
}
