package retry

type exec func ()(error, interface{});
type onError func(err error, count int)bool
func ToRetry(exec exec, onError onError)(err error, res interface{})  {
	i := 0
	for ;;i++ {
		err, res = exec()
		if err == nil{
			return
		}
		c := onError(err, i)
		if !c {
			return
		}
	}
}
