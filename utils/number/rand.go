package number

import (
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
)
func Rand(min int, max int) (num int, err error)  {
	var seed int64
	err = binary.Read(crand.Reader, binary.BigEndian, &seed)
	if err != nil {
		return 0, err
	}
	r := mrand.New(mrand.NewSource(seed))
	num = r.Intn(max - min) + min
	return
}
