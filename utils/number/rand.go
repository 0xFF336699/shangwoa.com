package number

import (
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"time"
)
func Random(min int, max int) (num int)  {
	mrand.Seed(int64(time.Now().UnixNano()))
	num = mrand.Intn(max - min) + min
	return
}
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


func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}