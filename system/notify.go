package system

import (
	"fmt"

	"shangwoa.com/log2"
)

func Notify(info string, level int, errs ...interface{}) {
	str := ""
	for _, e := range errs {
		if err, ok := e.(error); ok {
			str += err.Error()
		} else {
			str += fmt.Sprintf("%#v", e)
		}
	}
	log2.Errorln(info, str)
}
