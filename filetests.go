package main

import (
	"fmt"
	"shangwoa.com/os2"
)

func main() {
	err, isexists := os2.IsFile("xxx")
	fmt.Println(err, isexists)
}
