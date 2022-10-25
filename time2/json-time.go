package time2


import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// JsonTime format json time field by myself
type JsonTime struct {
	time.Time
}

// MarshalJSON on JsonTime format Time field with %Y-%m-%d %H:%M:%S
func (this *JsonTime) MarshalJSON() ([]byte, error) {

	formatted := fmt.Sprintf("\"%s\"", this.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (this *JsonTime) UnmarshalJSON(data []byte) (err error) {

	var t time.Time
	err, isOk, t := ParseTimeFromString(string(data))
	if err != nil{
		return
	}
	if isOk{
		*this = JsonTime{t}
	}
	return
}
// Value insert timestamp into mysql need this function.
func (t *JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Unix(t.Unix(), 10), nil
}

// Scan valueof time.Time
func (this *JsonTime) Scan(v interface{}) (err error) {
	if v == nil{
		return
	}
	var data []uint8
	switch v.(type) {
	case []uint8:
		data = v.([]uint8)
		break
	default:
		fmt.Println("unknow type", v)
		return errors.New("unknow type")
	}

	var t time.Time
	err, isOk, t := ParseTimeFromString(string(data))
	if err != nil{
		return
	}
	if isOk{
		*this = JsonTime{t}
	}
	return
}

