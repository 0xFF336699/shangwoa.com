package main

import (
	"encoding/base64"
	"shangwoa.com/crypto"
	"fmt"
)
func main() {
	result, err := crypto.AesEncrypt([]byte("hello hehe"), crypto.IvKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result))

	origData, err := crypto.AesDecrypt(result, crypto.IvKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
