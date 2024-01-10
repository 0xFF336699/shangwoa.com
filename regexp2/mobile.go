package regexp2

import (
	"regexp"
	"strings"
)

func IsMobile(mobile string) bool {
	re := regexp.MustCompile(`^1[0123456789][0-9]{9}$`)
	return re.MatchString(mobile)
}

func IsMobileOrMaskMobile(mobile string) bool {
	re := regexp.MustCompile(`^1[0123456789][0-9*]{9}$`)
	return re.MatchString(mobile)
}

func IsMaskMobile(mobile string) bool {
	bl := IsMobileOrMaskMobile(mobile)
	if !bl{
		return false
	}
	return strings.Index(mobile, "*") > -1
}

func hasMobileOrHasMaskMobile(mobile string) bool {
	re := regexp.MustCompile(`1[0123456789][0-9*]{9}`)
	return re.MatchString(mobile)
}
func HasMobile(mobile string) bool {
	re := regexp.MustCompile(`1[0123456789][0-9]{9}`)
	return re.MatchString(mobile)
}
