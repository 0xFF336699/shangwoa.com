package system

import (
	"runtime"
	"encoding/json"
	"fmt"
)
var DispatchError = false
var PrintError = false
const(
	SYSTEM_ERROR = "system_error"
)
type Error struct {
	Type      string
	ErrorFile string
	ErrorLine int
	Err       error
	SubType   string
	SubData   interface{}
	CodeLevel int
}

func (this *Error) Error() string {
	b, err := json.Marshal(this)
	if err != nil{
		return err.Error()
	}
	return string(b)
}
func OnError(err *Error)  {

	if DispatchError == false && PrintError == false{
		return
	}
	// 这个方法估计消耗会挺大，所以CodeLevel为0时就不调用了 节省一下资源，如果要调用就强制赋值不为0
	if err.CodeLevel > 0{
		_,file,line,_ := runtime.Caller(err.CodeLevel + 1)
		err.ErrorFile = file
		err.ErrorLine = line
	}
	if DispatchError{
		GetSystemDomainInstance().DispatchError(SYSTEM_ERROR, err)
		GetSystemDomainInstance().DispatchError(err.Type, err)
	}
	if PrintError{
		fmt.Println(err.Error())
	}
}