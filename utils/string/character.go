package string

import (
	"math/rand"
	"shangwoa.com/utils/number"
	"fmt"
)
const NumBytes = "0123456789"
const CharBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const CharNumBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const IV = "0123456789abcdef"

func RandString(letters string, n int) string {
	return string(RandStringBytes(letters, n))
}
func RandStringBytes(letters string, n int) ([]byte) {
	b := make([]byte, n)
	l := len(letters)
	for i := range b {
		index := rand.Intn(l)
		letter := letters[index]
		b[i] = letter
	}
	return b
}

func RandStringBytesWithRandLen(letters string, min int, max int) (b []byte, err error) {
	n, err := number.Rand(min, max)
	if err != nil {
		return nil, err
	}
	b = RandStringBytes(letters, n)
	return b, err
}

func RandStringWithRandLen(letters string, min int, max int) (str string, err error) {
	b, err := RandStringBytesWithRandLen(letters, min, max)
	str = string(b)
	fmt.Println("str is", str)
	return str, err
}