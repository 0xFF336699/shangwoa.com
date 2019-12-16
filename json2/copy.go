package json2

import (
	"encoding/base64"
	"encoding/json"
)

func Copy(in, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, out)
	return err
}

func ConvertPgArray(str string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	str = string(b)
	return str, nil
}
