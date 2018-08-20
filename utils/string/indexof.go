package string

func IndexOf(arr []string, item string) int {
	for k, v := range arr {
		if v == item{
			return  k
		}
	}
	return -1
}
