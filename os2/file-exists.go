package os2

import (
	"errors"
	"os"
)
var NotAFile = errors.New("it's not a file")
// 判断所给路径文件/文件夹是否存在
func Exists(path string) (error, bool) {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return nil, true
		}
		return nil, false
	}
	return nil,true
}

// 判断所给路径是否为文件夹
func IsDir(path string) (error, bool) {
	s, err := os.Stat(path)
	if err != nil {
		return err, false
	}
	return nil, s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) (error, bool) {
	err, bl := IsDir(path)
	return err, !bl
}