package time2

import (
	"time"
)

var HourDuration time.Duration
var MinuteDuration time.Duration

func init() {
	var err error
	HourDuration, err = time.ParseDuration("1h")
	if err != nil{
		panic(err)
	}
	MinuteDuration, err = time.ParseDuration("1m")
	if err != nil{
		panic(err)
	}
}

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

func ToFormatByLocation(in time.Time, location *time.Location, layout string) (string) {
	t := in.In(location)
	return t.Format(layout)
}

func TimeToBeiJingTime(in time.Time) ( t time.Time) {
	t = in.UTC()
	t = t.Add(8 * HourDuration)
	return
}

func AtTime(years, months, days , hours, minutes, seconds int) (t time.Time) {
	now := time.Now().UTC()
	now = now.AddDate(years, months, days)
	t = now.Add(HourDuration * time.Duration(hours) + MinuteDuration * time.Duration(minutes))
	return
}