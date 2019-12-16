package os2

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func Compare(spath, dpath string) (error, bool) {
	sinfo, err := os.Lstat(spath)
	if err != nil {
		fmt.Println("sinfo is", err.Error())
		return err, false
	}
	dinfo, err := os.Lstat(dpath)
	if err != nil {
		fmt.Println("dinfo is", err.Error())
		return err, false
	}
	if sinfo.Size() != dinfo.Size() {
		fmt.Println("Size is", sinfo.Size(), dinfo.Size())
		return err, false
	}
	return Comparefile(spath, dpath)
}

func Comparefile(spath, dpath string) (error, bool) {
	sFile, err := os.Open(spath)
	if err != nil {
		fmt.Println("sFile open is", err.Error())
		return err, false
	}
	dFile, err := os.Open(dpath)
	if err != nil {
		fmt.Println("dFile open is", err.Error())
		return err, false
	}
	err,b := Comparebyte(sFile, dFile)
	sFile.Close()
	dFile.Close()
	return err, b
}
//下面可以代替md5比较.
func Comparebyte(sfile *os.File, dfile *os.File) (error ,bool) {
	var sbyte []byte = make([]byte, 512)
	var dbyte []byte = make([]byte, 512)
	var serr, derr error
	for {
		_, serr = sfile.Read(sbyte)
		if serr != nil || derr != nil {
			if serr != derr {
				if serr != nil{
					fmt.Println("serr  is", serr.Error())
					return serr, false
				}else{
					fmt.Println("derr  is", derr.Error())
					return derr,false
				}

			}
			if serr == io.EOF {
				break
			}
		}
		if bytes.Equal(sbyte, dbyte) {
			continue
		}
		fmt.Println("compare equal is ", len(sbyte), len(dbyte))
		return nil, false
	}
	return nil, true
}
