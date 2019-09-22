package json2

import "encoding/json"

func MapToStruct(m interface{}, s interface{}) (err error)  {
	b, err := json.Marshal(m)
	if err != nil{return }
	err = json.Unmarshal(b, s)
	return
}
