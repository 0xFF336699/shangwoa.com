package regexp2

import "regexp"

func IsMobile(mobile string) bool {
	re := regexp.MustCompile(`^1[34578][0-9]{9}$`)
	return re.MatchString(mobile)
}

func HasMobile(mobile string) bool {
	re := regexp.MustCompile(`1[34578][0-9]{9}`)
	return re.MatchString(mobile)
}