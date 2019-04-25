package json2

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"

	"github.com/jmoiron/sqlx/types"
)

type JSONB map[string]interface{}

func (this *JSONB) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		if b[0] == []byte("[")[0] || b[0] == []byte("{")[0] {
			arr := map[string]interface{}{}
			if err := json.Unmarshal(b, &arr); err != nil {
				return err
			}
			*this = arr
			return nil
		} else {
			var s string
			if err := json.Unmarshal(b, &s); err != nil {
				return err
			}
			sb, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return err
			}
			err = this.Scan(sb)
			return err
		}
	}
	return nil
}
func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	if len(valueString) > 2 {
		str := string(valueString)
		if str == "null" {
			return "{}", err
		}
		return str, err
	}
	return "{}", err
}

func (j *JSONB) Scan(value interface{}) error {
	b := value.([]byte)
	if len(b) > 0 {
		m := map[string]interface{}{}
		if err := json.Unmarshal(value.([]byte), &m); err != nil {
			//fmt.Println("scan error", err)
			return err
		} // 为毛这样可以，下面就不行？
		*j = m
	} else {
		*j = nil
	}
	return nil
	//if err := json.Unmarshal(value.([]byte), &j); err != nil {
	//	return err
	//}
	//return nil
}

//func (j JSONB) Value() (driver.Value, error) {
//	valueString, err := json.Marshal(j)
//	return string(valueString), err
//}
//
//func (j *JSONB) Scan(value interface{}) error {
//	if err := json.Unmarshal(value.([]byte), &j); err != nil {
//		return err
//	}
//	return nil
//}

type JSONTags []string

func (tags *JSONTags) Scan(src interface{}) error {
	var jt types.JSONText

	if err := jt.Scan(src); err != nil {
		return err
	}

	if err := jt.Unmarshal(tags); err != nil {
		return err
	}

	return nil
}

func (tags *JSONTags) Value() (driver.Value, error) {
	var jt types.JSONText

	data, err := json.Marshal((*[]string)(tags))
	if err != nil {
		return nil, err
	}

	if err := jt.UnmarshalJSON(data); err != nil {
		return nil, err
	}

	return jt.Value()
}

func (tags *JSONTags) MarshalJSON() ([]byte, error) {
	return json.Marshal((*[]string)(tags))
}

func (tags *JSONTags) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*[]string)(tags)); err != nil {
		return err
	}

	return nil
}

type ArrayJsonb []JSONB

func (a ArrayJsonb) Value() (driver.Value, error) {

	valueString, err := json.Marshal(a)
	if len(valueString) > 2 {
		str := string(valueString)
		if str == "null" {
			return "{}", err
		}
		return str, err
	}
	return "{}", err
}

func (j *ArrayJsonb) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

func (this *ArrayJsonb) UnmarshalJSON(b []byte) error {
	if len(b) > 0 {
		if b[0] == []byte("[")[0] {
			arr := []JSONB{}
			if err := json.Unmarshal(b, &arr); err != nil {
				return err
			}
			*this = arr
			return nil
		} else {
			var s string
			if err := json.Unmarshal(b, &s); err != nil {
				return err
			}
			sb, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return err
			}
			err = this.Scan(sb)
			return err
		}
	}
	return nil
}

//
//type ArrayJsonb []JSONB
//
//func (this *ArrayJsonb) Scan(src interface{}) error {
//	return nil
//}
//
//func (this ArrayJsonb) Value() (driver.Value, error) {
//	b, err:= json.Marshal(this)
//	if len(b) > 2{
//		b[0] = []byte("{")[0]
//		b[len(b) -1] = []byte("}")[0]
//		//b = b[1:len(b) - 1]
//		str := string(b)
//		fmt.Println(str)
//		return string(b), err
//	}
//	return "{}", err
//}
//
func TransformType(src interface{}) (interface{}, error) {
	switch src.(type) {
	case []byte:
		b := src.([]byte)
		j := &JSONB{}
		err := json.Unmarshal(b, j)
		if err != nil {
			return src, err
		}
		return j, nil
	}
	return src, nil
}
