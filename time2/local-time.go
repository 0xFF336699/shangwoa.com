package time2

import (
	"time"
)

func ToLocation(in *time.Time, location string)(err error, out *time.Time)  {
	l, err := time.LoadLocation(location)
	if err != nil{return }
	t := in.In(l)
	out = &t
	return
}

func ToLocal(in *time.Time) (err error, out *time.Time) {
	return ToLocation(in, "Local")
}

func ToUTC(in *time.Time) (err error, out *time.Time) {
	return ToLocation(in, "UTC")
}

func GetUCT() (err error, out *time.Time) {
	now := time.Now()
	return ToUTC(&now)
}
func ToLocalDefaultFormat(in *time.Time) (err error, str string) {
	err, t := ToLocal(in)
	if err != nil {return }
	str = t.Format(FormatOnlyTime)
	return
}