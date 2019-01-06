package main

import (
	"regexp"
	"strconv"
	"strings"
)

func main() {
	test2()
}
func test1()  {
	message := `盗图 4`
	reg := regexp.MustCompile(`(^盗图|^盗图狗|^搬运表情|^表情搬运|^偷盗表情|^偷表情|^订阅表情)(\d*|[ ]*|[ ]*\d*)$`)
	res := reg.FindStringSubmatch(message)
	if len(res) == 3{
		s := strings.Replace(res[2], " ", "", -1)
		n, err := strconv.Atoi(s)
		if err != nil{
			n = 0
		}
		max := 10
		if n == 0 || n < max{
			n = max
		}
		println(n)
	}
	println(len(res), res)
}

func test2() {
	message := "fc2"
	reg := regexp.MustCompile(`(^fc2|^来一波日本妹子|^来一波岛国妹子|^岛国妹子来一波|^日本妹子来一波)([ ]*)$`)
	res := reg.FindStringSubmatch(message)
	if len(res) == 3{
		return
	}
}