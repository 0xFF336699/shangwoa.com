package number

import (
	"errors"
	"strconv"
)

func InterfaceToInt(val interface{}, forceString bool) (i int, err error) {
	switch t := val.(type) {
	case int:
		i = t
	case int8:
		i = int(t) // standardizes across systems
	case int16:
		i = int(t) // standardizes across systems
	case int32:
		i = int(t) // standardizes across systems
	case int64:
		i = int(t) // standardizes across systems
	case bool:
		err = errors.New("is bool")
		// // not covertible unless...
		// if t {
		//  i = 1
		// } else {
		//  i = 0
		// }
	case float32:
		i = int(t) // standardizes across systems
	case float64:
		i = int(t) // standardizes across systems
	case uint8:
		i = int(t) // standardizes across systems
	case uint16:
		i = int(t) // standardizes across systems
	case uint32:
		i = int(t) // standardizes across systems
	case uint64:
		i = int(t) // standardizes across systems
	case string:
		if forceString{
			s := val.(string)
			i, err= strconv.Atoi(s)
		}
		// gets a little messy...
	default:
		// what is it then?
		err = errors.New("not match type")
	}
	return
}
