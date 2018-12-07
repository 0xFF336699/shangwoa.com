package array

import (
	"reflect"
	"sort"
	"testing"
)
type RGB struct {
	R, G, B uint8
}
type BGR struct {
	B, G, R uint8
}

func RGB2BGR(data []RGB) []BGR {
	d := Slice(data, reflect.TypeOf([]BGR(nil)))
	return d.([]BGR)
}
func TestIndexOf(t *testing.T) {
	i := int(132103472)
	println("i is", i)
	return
	l := []string{"a", "b"}
	needle := "b"
	idx := sort.Search(len(l), func(i int) bool {
		return l[i] >= needle
	})
	println(idx)
}

func IndexOf(params ...interface{}) int {
	v := reflect.ValueOf(params[0])
	arr := reflect.ValueOf(params[1])

	var t = reflect.TypeOf(params[1]).Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len()-1; i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return i
		}
	}
	return -1
}
